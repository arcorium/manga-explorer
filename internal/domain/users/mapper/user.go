package mapper

import (
	"manga-explorer/internal/domain/users"
	"manga-explorer/internal/domain/users/dto"
	"manga-explorer/internal/util"
	"time"
)

func ToUserResponse(user *users.User) dto.UserResponse {
	return dto.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.Underlying(),
	}
}

// MapUserRegisterInput Create new user and the password is hashed automatically and email will be validated, it will return error ErrEmailValidation or ErrHashPassword or nil
func MapUserRegisterInput(input *dto.UserRegisterInput) (users.User, error) {
	return users.NewUser(input.Username, input.Email, input.Password, users.RoleUser)
}

// MapUserUpdateInput Create user for update except password
func MapUserUpdateInput(input *dto.UpdateUserInput) users.User {
	return users.User{
		Id:        input.UserId,
		Username:  input.Username,
		Email:     input.Email,
		UpdatedAt: time.Now(),
	}
}

func createUserForPasswordChange(userId, password string) (users.User, error) {
	usr := users.User{
		Id:        userId,
		UpdatedAt: time.Now(),
	}

	passwords, err := util.Hash(password)
	if err != nil {
		return usr, users.ErrHashPassword
	}
	usr.Password = passwords
	return usr, nil
}

// MapChangePasswordInput Create user for update password, the password will be hashed automatically
func MapChangePasswordInput(input *dto.ChangePasswordInput) (users.User, error) {
	return createUserForPasswordChange(input.UserId, input.NewPassword)
}

func MapResetPasswordInput(input *dto.ResetPasswordInput) (users.User, error) {
	return createUserForPasswordChange(input.UserId, input.NewPassword)
}

func MapAddUserInput(input *dto.AddUserInput) (users.User, error) {
	role := users.Role(input.Role)
	if err := role.Validate(); err != nil {
		return users.User{}, err
	}
	return users.NewUser(input.Username, input.Email, input.Password, role)
}

func MapUpdateUserExtendedInput(input *dto.UpdateUserExtendedInput) (users.User, error) {
	user := users.User{
		Id:        input.UserId,
		Username:  input.Username,
		Email:     input.Email,
		UpdatedAt: time.Now(),
	}

	if len(input.Password) == 0 {
		return user, nil
	}

	password, err := util.Hash(input.Password)
	if err != nil {
		return user, err
	}
	user.Password = password

	return user, nil
}