package service

import (
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/domain/users/dto"
)

type IUser interface {
	// RegisterUser creating user and profile belong to it based on the input
	RegisterUser(request *dto.UserRegisterInput) (dto.UserResponse, common.Status)
	// AddUser creating user and profile for all possible field based on the input
	AddUser(input *dto.AddUserInput) common.Status
	// DeleteUser delete the user based on the id
	DeleteUser(userId string) common.Status
	// GetAllUsers return all the users
	GetAllUsers() ([]dto.UserResponse, common.Status)
	// FindUserByEmail find user based on the email
	FindUserByEmail(email string) (dto.UserResponse, common.Status)
	// FindUserById find user by the user id
	FindUserById(id string) (dto.UserResponse, common.Status)
	// FindUserProfileById Get user and the profile
	FindUserProfileById(userId string) (dto.ProfileResponse, common.Status)
	// FindMangaHistories find all user manga histories

	// UpdateUser update user (not profile) based on the input except for password field which should be handled by ChangePassword
	UpdateUser(input *dto.UpdateUserInput) common.Status
	// UpdateUserExtended update user including the password, the method should handle hashing the password before store it into the persistent storage
	UpdateUserExtended(input *dto.UpdateUserExtendedInput) common.Status
	// UpdateProfile update user profile based on user id (not profile id) and the input
	UpdateProfile(input *dto.ProfileUpdateInput) common.Status
	// UpdateProfileExtended update all possible field on user profile
	UpdateProfileExtended(input *dto.UpdateProfileExtendedInput) common.Status
	// ChangePassword works like UpdateUser but it specific for password
	ChangePassword(input *dto.ChangePasswordInput) common.Status
	// ResetPassword works like ChangePassword, but it will not validate the last password, instead it will just patch the password
	ResetPassword(input *dto.ResetPasswordInput) common.Status
}
