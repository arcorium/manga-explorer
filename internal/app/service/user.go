package service

import (
	"manga-explorer/internal/app/common/constant"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/mapper"
	"manga-explorer/internal/domain/users/repository"
	"manga-explorer/internal/domain/users/service"
	"manga-explorer/internal/infrastructure/file"
	fileService "manga-explorer/internal/infrastructure/file/service"
	"manga-explorer/internal/infrastructure/mail"
	mailService "manga-explorer/internal/infrastructure/mail/service"
	"manga-explorer/internal/util/containers"
	"time"
)

func NewUser(userRepo repository.IUser, verification service.IVerification, authentication service.IAuthentication, mail mailService.IMail) service.IUser {
	return &userService{repo: userRepo, verifService: verification, mailService: mail, authService: authentication}
}

type userService struct {
	repo         repository.IUser
	verifService service.IVerification
	authService  service.IAuthentication
	mailService  mailService.IMail
	fileService  fileService.IFile
}

func (u userService) UpdateProfileImage(input *dto.ProfileImageUpdateInput) status.Object {
	// Get user profile
	profiles, err := u.repo.FindUserProfiles(input.UserId)
	if err != nil {
		return status.RepositoryError(err)
	}

	// Upload new image
	filename, stat := u.fileService.Upload(file.ProfileAsset, input.Image)
	if stat.IsError() {
		return stat
	}

	// Delete old image
	if len(profiles.PhotoURL) > 0 {
		stat = u.fileService.Delete(file.ProfileAsset, profiles.PhotoURL)
		if stat.IsError() {
			return stat
		}
	}

	// Set metadata
	profile := users.Profile{Id: profiles.Id, PhotoURL: filename, UpdatedAt: time.Now()}
	err = u.repo.UpdateProfile(&profile)
	return status.ConditionalRepository(err, status.UPDATED)
}

func (u userService) DeleteProfileImage(userId string) status.Object {
	profiles, err := u.repo.FindUserProfiles(userId)
	if err != nil {
		return status.RepositoryError(err)
	}
	stat := u.fileService.Delete(file.ProfileAsset, profiles.PhotoURL)
	if stat.IsError() {
		return stat
	}

	profile := users.Profile{Id: profiles.Id, PhotoURL: "", UpdatedAt: time.Now()}
	err = u.repo.UpdateProfile(&profile)
	return status.ConditionalRepository(err, status.DELETED)
}

func (u userService) AddUser(input *dto.AddUserInput) status.Object {
	userInput, err := mapper.MapAddUserInput(input)
	if err != nil {
		return status.InternalError()
	}
	profileInput := mapper.MapAddProfileInput(&userInput, input)

	err = u.repo.CreateUser(&userInput, &profileInput)
	return status.ConditionalRepository(err, status.CREATED)
}
func (u userService) DeleteUser(userId string) status.Object {
	err := u.repo.DeleteUser(userId)
	return status.ConditionalRepository(err, status.DELETED)
}
func (u userService) GetAllUsers() ([]dto.UserResponse, status.Object) {
	allUsers, err := u.repo.GetAllUsers()
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	result := containers.CastSlicePtr(allUsers, mapper.ToUserResponse)
	return result, status.Success()
}

func (u userService) UpdateUser(input *dto.UpdateUserInput) status.Object {
	usr := mapper.MapUserUpdateInput(input)
	err := u.repo.UpdateUser(&usr)
	return status.ConditionalRepository(err, status.UPDATED)
}
func (u userService) UpdateUserExtended(input *dto.UpdateUserExtendedInput) status.Object {
	user, err := mapper.MapUpdateUserExtendedInput(input)
	if err != nil {
		return status.InternalError()
	}
	err = u.repo.UpdateUser(&user)
	return status.ConditionalRepository(err, status.UPDATED)
}
func (u userService) UpdateProfileExtended(input *dto.UpdateProfileExtendedInput) status.Object {
	profile := mapper.MapUpdateProfileExtendedInput(input)
	err := u.repo.UpdateProfile(&profile)
	return status.ConditionalRepository(err, status.UPDATED)
}
func (u userService) UpdateProfile(input *dto.ProfileUpdateInput) status.Object {
	profile := mapper.MapProfileUpdateInput(input)
	err := u.repo.UpdateProfile(&profile)
	return status.ConditionalRepository(err, status.UPDATED)
}
func (u userService) RegisterUser(input *dto.UserRegisterInput) (dto.UserResponse, status.Object) {
	usr, err := mapper.MapUserRegisterInput(input)
	if err != nil {
		return dto.UserResponse{}, status.InternalError()
	}

	profile := mapper.MapProfileRegisterInput(&usr, input)
	err = u.repo.CreateUser(&usr, &profile)
	if err != nil {
		return dto.UserResponse{}, status.RepositoryError(err)
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
	if err != nil {
		return dto.UserResponse{}, status.RepositoryError(err)
	}
	return mapper.ToUserResponse(usr), status.Success()
}
func (u userService) FindUserByEmail(email string) (dto.UserResponse, status.Object) {
	usr, err := u.repo.FindUserByEmail(email)
	if err != nil {
		return dto.UserResponse{}, status.RepositoryError(err)
	}
	return mapper.ToUserResponse(usr), status.Success()
}
func (u userService) FindUserProfileById(userId string) (dto.ProfileResponse, status.Object) {
	profile, err := u.repo.FindUserProfiles(userId)
	if err != nil {
		return dto.ProfileResponse{}, status.RepositoryError(err)
	}
	return mapper.ToProfileResponse(profile, u.fileService), status.Success()
}
func (u userService) ChangePassword(input *dto.ChangePasswordInput) status.Object {
	// Get user
	usr, err := u.repo.FindUserById(input.UserId)
	if err != nil {
		return status.RepositoryError(err)
	}

	// Check last password
	if !usr.ValidatePassword(input.LastPassword) {
		return status.Error(status.USER_LOGIN_ERROR)
	}

	// Set new password
	updateUser, err := mapper.MapChangePasswordInput(input)
	if err != nil {
		return status.InternalError()
	}
	err = u.repo.UpdateUser(&updateUser)
	return status.ConditionalRepository(err, status.UPDATED)
}
func (u userService) RequestResetPassword(input *dto.ResetPasswordRequestInput) status.Object {
	userResponse, err := u.repo.FindUserByEmail(input.Email)
	if err != nil {
		return status.RepositoryError(err)
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
	updateUser, err := mapper.MapResetPasswordInput(input, verifResponse.UserId)
	if err != nil {
		return status.InternalError()
	}
	err = u.repo.UpdateUser(&updateUser)
	return status.ConditionalRepository(err, status.SUCCESS)
}
