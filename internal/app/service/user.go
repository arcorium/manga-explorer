package service

import (
	"manga-explorer/internal/app/common/constant"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/mapper"
	"manga-explorer/internal/domain/users/repository"
	"manga-explorer/internal/domain/users/service"
	"manga-explorer/internal/infrastructure/mail"
	mailService "manga-explorer/internal/infrastructure/mail/service"
	"manga-explorer/internal/util/containers"
)

func NewUser(userRepo repository.IUser, verification service.IVerification, authentication service.IAuthentication, mail mailService.IMail) service.IUser {
	return &userService{repo: userRepo, verifService: verification, mailService: mail, authService: authentication}
}

type userService struct {
	repo         repository.IUser
	verifService service.IVerification
	mailService  mailService.IMail
	authService  service.IAuthentication
}

func (u userService) AddUser(input *dto.AddUserInput) status.Object {
	userInput, err := mapper.MapAddUserInput(input)
	if err != nil {
		return status.Error(status.INTERNAL_SERVER_ERROR)
	}
	profileInput := mapper.MapAddProfileInput(&userInput, input)

	return status.FromRepository(u.repo.CreateUser(&userInput, &profileInput), status.CREATED)
}

func (u userService) DeleteUser(userId string) status.Object {
	return status.FromRepository(u.repo.DeleteUser(userId), status.DELETED)
}

func (u userService) GetAllUsers() ([]dto.UserResponse, status.Object) {
	allUsers, err := u.repo.GetAllUsers()
	return containers.CastSlicePtr(allUsers, mapper.ToUserResponse), status.FromRepository(err, status.SUCCESS)
}

func (u userService) UpdateUser(input *dto.UpdateUserInput) status.Object {
	usr := mapper.MapUserUpdateInput(input)
	err := u.repo.UpdateUser(&usr)
	return status.FromRepository(err, status.UPDATED)
}

func (u userService) UpdateUserExtended(input *dto.UpdateUserExtendedInput) status.Object {
	user, err := mapper.MapUpdateUserExtendedInput(input)
	if err != nil {
		return status.Error(status.INTERNAL_SERVER_ERROR)
	}
	return status.FromRepository(u.repo.UpdateUser(&user), status.UPDATED)
}

func (u userService) UpdateProfileExtended(input *dto.UpdateProfileExtendedInput) status.Object {
	profile := mapper.MapUpdateProfileExtendedInput(input)
	return status.FromRepository(u.repo.UpdateProfile(&profile), status.UPDATED)
}

func (u userService) UpdateProfile(input *dto.ProfileUpdateInput) status.Object {
	profile := mapper.MapProfileUpdateInput(input)
	err := u.repo.UpdateProfile(&profile)
	return status.FromRepository(err, status.UPDATED)
}
func (u userService) RegisterUser(input *dto.UserRegisterInput) (dto.UserResponse, status.Object) {
	usr, err := mapper.MapUserRegisterInput(input)
	if err != nil {
		return dto.UserResponse{}, status.Error(status.INTERNAL_SERVER_ERROR)
	}

	profile := mapper.MapProfileRegisterInput(&usr, input)
	err = u.repo.CreateUser(&usr, &profile)
	stat := status.FromRepository(err, status.CREATED)
	if stat.IsError() {
		return dto.UserResponse{}, stat
	}

	// Send email Verification
	verifResponse, stat := u.verifService.Request(usr.Id, users.UsageVerifyEmail)
	if stat.IsError() {
		return dto.UserResponse{}, stat
	}

	ml := &mail.Mail{
		From:    constant.IssuerName,
		To:      usr.Email,
		Subject: "Email Verification",
		Body:    verifResponse.Token,
	}
	go u.mailService.SendEmail(ml)

	return mapper.ToUserResponse(&usr), stat
}
func (u userService) FindUserById(id string) (dto.UserResponse, status.Object) {
	usr, err := u.repo.FindUserById(id)
	stat := status.FromRepository(err)
	if stat.IsError() {
		return dto.UserResponse{}, stat
	}
	return mapper.ToUserResponse(usr), stat
}
func (u userService) FindUserByEmail(email string) (dto.UserResponse, status.Object) {
	usr, err := u.repo.FindUserByEmail(email)
	stat := status.FromRepository(err)
	if stat.IsError() {
		return dto.UserResponse{}, stat
	}
	return mapper.ToUserResponse(usr), stat
}
func (u userService) FindUserProfileById(userId string) (dto.ProfileResponse, status.Object) {
	profile, err := u.repo.FindUserProfiles(userId)
	stat := status.FromRepository(err)
	if stat.IsError() {
		return dto.ProfileResponse{}, stat
	}
	return mapper.ToProfileResponse(profile), stat
}

func (u userService) ChangePassword(input *dto.ChangePasswordInput) status.Object {
	// Get user
	usr, err := u.repo.FindUserById(input.UserId)
	if err != nil {
		return status.Error(status.USER_NOT_FOUND)
	}

	// Check last password
	if !usr.ValidatePassword(input.LastPassword) {
		return status.Error(status.USER_LOGIN_ERROR)
	}

	// Set new password
	updateUser, err := mapper.MapChangePasswordInput(input)
	if err != nil {
		return status.Error(status.INTERNAL_SERVER_ERROR)
	}
	err = u.repo.UpdateUser(&updateUser)
	return status.FromRepository(err, status.UPDATED)
}

func (u userService) RequestResetPassword(input *dto.ResetPasswordRequestInput) status.Object {
	userResponse, err := u.repo.FindUserByEmail(input.Email)
	stat := status.FromRepository(err)
	if stat.IsError() {
		return stat
	}

	verifResponse, stat := u.verifService.Request(userResponse.Id, users.UsageResetPassword)
	if stat.IsError() {
		return stat
	}

	ml := &mail.Mail{
		From:    constant.IssuerName,
		To:      userResponse.Email,
		Subject: "Reset Password",
		Body:    verifResponse.Token,
	}
	go u.mailService.SendEmail(ml)
	return status.Success()
}

func (u userService) ResetPassword(input *dto.ResetPasswordInput) status.Object {
	// Find token
	verifResponse, stat := u.verifService.Find(input.Token)
	if stat.IsError() {
		return stat
	}
	input.UserId = verifResponse.UserId

	// Validate token
	stat = u.verifService.Validate(&verifResponse)
	if stat.IsError() {
		return stat
	}

	// Immediately revoke verify token
	stat = u.verifService.Remove(input.Token)
	if stat.IsError() {
		return stat
	}

	// Logout all devices
	u.authService.LogoutDevices(verifResponse.UserId)

	// Set new password
	updateUser, err := mapper.MapResetPasswordInput(input)
	if err != nil {
		return status.Error(status.INTERNAL_SERVER_ERROR)
	}
	err = u.repo.UpdateUser(&updateUser)
	return status.FromRepository(err, status.UPDATED)
}
