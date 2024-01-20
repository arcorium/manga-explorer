package service

import (
	"errors"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	mangaServices "manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/domain/users/mapper"
	"manga-explorer/internal/domain/users/repository"
	"manga-explorer/internal/domain/users/service"
	"manga-explorer/internal/util/containers"
)

func NewUser(userRepo repository.IUser, mangaService mangaServices.IManga) service.IUser {
	return &userService{userRepo: userRepo, mangaService: mangaService}
}

type userService struct {
	userRepo     repository.IUser
	mangaService mangaServices.IManga
}

func (u userService) AddUser(input *dto.AddUserInput) common.Status {
	userInput, err := mapper.MapAddUserInput(input)
	if err != nil {
		return common.StatusError(status.INTERNAL_SERVER_ERROR)
	}
	profileInput := mapper.MapAddProfileInput(&userInput, input)

	return common.NewRepositoryStatus(u.userRepo.CreateUser(&userInput, &profileInput), status.SUCCESS_CREATED)
}

func (u userService) DeleteUser(userId string) common.Status {
	return common.NewRepositoryStatus(u.userRepo.DeleteUser(userId))
}

func (u userService) GetAllUsers() ([]dto.UserResponse, common.Status) {
	allUsers, err := u.userRepo.GetAllUsers()
	return containers.CastSlicePtr(allUsers, mapper.ToUserResponse), common.NewRepositoryStatus(err)
}

func (u userService) UpdateUserExtended(input *dto.UpdateUserExtendedInput) common.Status {
	user, err := mapper.MapUpdateUserExtendedInput(input)
	if err != nil {
		return common.StatusError(status.INTERNAL_SERVER_ERROR)
	}

	return common.NewRepositoryStatus(u.userRepo.UpdateUser(&user))
}

func (u userService) UpdateProfileExtended(input *dto.UpdateProfileExtendedInput) common.Status {
	profile := mapper.MapUpdateProfileExtendedInput(input)
	return common.NewRepositoryStatus(u.userRepo.UpdateProfile(&profile))
}

func (u userService) UpdateProfile(input *dto.ProfileUpdateInput) common.Status {
	profile := mapper.MapProfileUpdateInput(input)
	err := u.userRepo.UpdateProfile(&profile)
	cerr := common.NewRepositoryStatus(err)
	return cerr
}
func (u userService) RegisterUser(input *dto.UserRegisterInput) (dto.UserResponse, common.Status) {
	usr, err := mapper.MapUserRegisterInput(input)
	if err != nil {
		if errors.Is(err, users.ErrEmailValidation) {
			return dto.UserResponse{}, common.StatusError(status.USER_CREATION_BAD_EMAIL)
		}

		return dto.UserResponse{}, common.StatusError(status.INTERNAL_SERVER_ERROR)
	}

	profile := mapper.MapProfileRegisterInput(&usr, input)
	err = u.userRepo.CreateUser(&usr, &profile)
	status := common.NewRepositoryStatus(err)
	if status.IsError() {
		return dto.UserResponse{}, status
	}
	return mapper.ToUserResponse(&usr), common.StatusSuccess()
}
func (u userService) FindUserById(id string) (dto.UserResponse, common.Status) {
	usr, err := u.userRepo.FindUserById(id)
	if err != nil {
		return dto.UserResponse{}, common.StatusError(status.INTERNAL_SERVER_ERROR)
	}
	return mapper.ToUserResponse(usr), common.StatusSuccess()
}
func (u userService) FindUserByEmail(email string) (dto.UserResponse, common.Status) {
	usr, err := u.userRepo.FindUserByEmail(email)
	if err != nil {
		return dto.UserResponse{}, common.StatusError(status.INTERNAL_SERVER_ERROR)
	}
	return mapper.ToUserResponse(usr), common.StatusSuccess()
}
func (u userService) FindUserProfileById(userId string) (dto.ProfileResponse, common.Status) {
	profile, err := u.userRepo.FindUserProfiles(userId)
	if err != nil {
		return dto.ProfileResponse{}, common.StatusError(status.INTERNAL_SERVER_ERROR)
	}
	return mapper.ToProfileResponse(profile), common.StatusSuccess()
}
func (u userService) UpdateUser(input *dto.UpdateUserInput) common.Status {
	usr := mapper.MapUserUpdateInput(input)
	err := u.userRepo.UpdateUser(&usr)
	if err != nil {
		return common.StatusError(status.INTERNAL_SERVER_ERROR)
	}
	return common.StatusSuccess()
}
func (u userService) ChangePassword(input *dto.ChangePasswordInput) common.Status {
	// Get user
	usr, err := u.userRepo.FindUserById(input.UserId)
	if err != nil {
		return common.StatusError(status.USER_NOT_FOUND)
	}

	// Check last password
	if !usr.ValidatePassword(input.LastPassword) {
		return common.StatusError(status.USER_LOGIN_ERROR)
	}

	// Set new password
	updateUser, err := mapper.MapChangePasswordInput(input)
	if err != nil {
		return common.StatusError(status.INTERNAL_SERVER_ERROR)
	}
	err = u.userRepo.UpdateUser(&updateUser)
	return common.ConditionalStatus(err, status.USER_UPDATE_FAILED)
}
func (u userService) ResetPassword(input *dto.ResetPasswordInput) common.Status {
	// Check user
	_, err := u.userRepo.FindUserById(input.UserId)
	if err != nil {
		return common.StatusError(status.USER_NOT_FOUND)
	}

	// Set new password
	updateUser, err := mapper.MapResetPasswordInput(input)
	if err != nil {
		return common.StatusError(status.INTERNAL_SERVER_ERROR)
	}
	err = u.userRepo.UpdateUser(&updateUser)
	return common.ConditionalStatus(err, status.USER_UPDATE_FAILED)
}
