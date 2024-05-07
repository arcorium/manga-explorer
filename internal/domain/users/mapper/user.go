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
    Role:     user.Role.String(),
  }
}

// MapUserRegisterInput Upsert new user and the password is hashed automatically and email will be validated, it will return error ErrEmailValidation or ErrHashPassword or nil
func MapUserRegisterInput(input *dto.UserRegisterInput) (users.User, error) {
  return users.NewUser(input.Username, input.Email, input.Password, users.RoleUser)
}

// MapUserUpdateInput Upsert user for update except password
func MapUserUpdateInput(input *dto.UserEditInput) users.User {
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
    return users.BadUser, users.ErrHashPassword
  }
  usr.Password = passwords
  return usr, nil
}

// MapChangePasswordInput Upsert user for update password, the password will be hashed automatically
func MapChangePasswordInput(input *dto.ChangePasswordInput) (users.User, error) {
  return createUserForPasswordChange(input.UserId, input.NewPassword)
}

func MapResetPasswordInput(input *dto.ResetPasswordInput, userId string) (users.User, error) {
  return createUserForPasswordChange(userId, input.NewPassword)
}

func MapVerifyEmailInput(userId string) users.User {
  return users.User{
    Id:        userId,
    Verified:  true,
    UpdatedAt: time.Now(),
  }
}

func MapAddUserInput(input *dto.AddUserInput) (users.User, error) {
  role := users.NewRole(input.Role)
  if err := role.Validate(); err != nil {
    return users.BadUser, err
  }
  return users.NewUser(input.Username, input.Email, input.Password, role)
}

func MapUserUpdateExtendedInput(input *dto.UserEditExtendedInput) (users.User, error) {
  user := users.User{
    Id:        input.UserId,
    Username:  input.Username,
    Email:     input.Email,
    UpdatedAt: time.Now(),
  }

  if len(input.Password) == 0 {
    return users.BadUser, users.ErrHashPassword
  }

  password, err := util.Hash(input.Password)
  if err != nil {
    return users.BadUser, err
  }
  user.Password = password

  return user, nil
}
