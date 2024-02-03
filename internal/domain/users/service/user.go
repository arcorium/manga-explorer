package service

import (
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/users/dto"
)

type IUser interface {
	// RegisterUser creating user and profile belong to it based on the input
	RegisterUser(request *dto.UserRegisterInput) (dto.UserResponse, status.Object)
	// AddUser creating user and profile for all possible field based on the input
	AddUser(input *dto.AddUserInput) status.Object
	// DeleteUser delete the user based on the id
	DeleteUser(userId string) status.Object
	// GetAllUsers return all the users
	GetAllUsers() ([]dto.UserResponse, status.Object)
	// FindUserByEmail find user based on the email
	FindUserByEmail(email string) (dto.UserResponse, status.Object)
	// FindUserById find user by the user id
	FindUserById(id string) (dto.UserResponse, status.Object)
	// FindUserProfileById Get user and the profile
	FindUserProfileById(userId string) (dto.ProfileResponse, status.Object)
	// UpdateUser update user (not profile) based on the input except for password field which should be handled by ChangePassword
	UpdateUser(input *dto.UserEditInput) status.Object
	// UpdateUserExtended update user including the password, the method should handle hashing the password before store it into the persistent storage
	UpdateUserExtended(input *dto.UserEditExtendedInput) status.Object
	// UpdateProfile update user profile based on user id (not profile id) and the input
	UpdateProfile(input *dto.ProfileEditInput) status.Object
	UpdateProfileImage(input *dto.ProfileImageUpdateInput) status.Object
	DeleteProfileImage(userId string) status.Object
	// UpdateProfileExtended update all possible field on user profile
	UpdateProfileExtended(input *dto.ProfileEditExtendedInput) status.Object
	// ChangePassword works like UpdateUser but it specific for password
	ChangePassword(input *dto.ChangePasswordInput) status.Object
	RequestResetPassword(input *dto.ResetPasswordRequestInput) status.Object
	// ResetPassword works like ChangePassword, but it will not validate the last password, instead it will just patch the password
	ResetPassword(input *dto.ResetPasswordInput) status.Object
	RequestEmailVerification(input *dto.VerifEmailRequestInput) status.Object
	VerifyEmail(input *dto.VerifyEmailInput) status.Object
}
