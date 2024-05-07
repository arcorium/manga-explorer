package service

import (
  "fmt"
  "log"
  "manga-explorer/internal/common"
  "manga-explorer/internal/common/status"
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
  "manga-explorer/internal/util/opt"
  "time"
)

func NewUser(config *common.Config, userRepo repository.IUser, verification service.IVerification, authentication service.IAuthentication, mail mailService.IMail, file fileService.IFile) service.IUser {
  return &userService{config: config, repo: userRepo, verifService: verification, mailService: mail, fileService: file, authService: authentication}
}

type userService struct {
  config *common.Config

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
    return status.RepositoryError(err, opt.New(status.PROFILE_NOT_FOUND))
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
      u.fileService.Delete(file.ProfileAsset, filename) // Delete uploaded image due to deletion error
      return stat
    }
  }

  // Set metadata
  profile := users.Profile{Id: profiles.Id, PhotoURL: filename, UpdatedAt: time.Now()}
  err = u.repo.UpdateProfile(&profile)
  if err != nil {
    u.fileService.Delete(file.ProfileAsset, filename)
    return status.RepositoryError(err, opt.New(status.PROFILE_UPDATE_FAILED))
  }
  return status.Updated()
}

func (u userService) DeleteProfileImage(userId string) status.Object {
  profiles, err := u.repo.FindUserProfiles(userId)
  if err != nil {
    return status.RepositoryError(err, opt.New(status.PROFILE_NOT_FOUND))
  }

  // Do nothing when the profile has no image
  if len(profiles.PhotoURL) == 0 {
    return status.Deleted()
  }

  stat := u.fileService.Delete(file.ProfileAsset, profiles.PhotoURL)
  if stat.IsError() {
    return stat
  }

  profile := users.Profile{Id: profiles.Id, PhotoURL: file.NoFile, UpdatedAt: time.Now()}
  err = u.repo.UpdateProfile(&profile)
  return status.ConditionalRepository(err, status.DELETED, opt.New(status.PROFILE_UPDATE_FAILED))
}

func (u userService) AddUser(input *dto.AddUserInput) status.Object {
  userInput, err := mapper.MapAddUserInput(input)
  if err != nil {
    return status.BadRequestError()
  }
  profileInput := mapper.MapAddProfileInput(&userInput, input)

  err = u.repo.CreateUser(&userInput, &profileInput)
  return status.ConditionalRepositoryE(err, status.CREATED, opt.New(status.USER_CREATION_ALREADY_EXIST), opt.New(status.USER_CREATION_ALREADY_EXIST))
}

func (u userService) DeleteUser(userId string) status.Object {
  err := u.repo.DeleteUser(userId)
  return status.ConditionalRepository(err, status.DELETED, opt.New(status.USER_NOT_FOUND))
}

func (u userService) GetAllUsers() ([]dto.UserResponse, status.Object) {
  allUsers, err := u.repo.GetAllUsers()
  result := containers.CastSlicePtr(allUsers, mapper.ToUserResponse)
  return result, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (u userService) UpdateUser(input *dto.UserEditInput) status.Object {
  usr := mapper.MapUserUpdateInput(input)
  err := u.repo.UpdateUser(&usr)
  return status.ConditionalRepository(err, status.UPDATED, opt.New(status.USER_UPDATE_FAILED))
}

func (u userService) UpdateUserExtended(input *dto.UserEditExtendedInput) status.Object {
  user, err := mapper.MapUserUpdateExtendedInput(input)
  if err != nil {
    return status.InternalError()
  }
  err = u.repo.UpdateUser(&user)
  return status.ConditionalRepository(err, status.UPDATED, opt.New(status.USER_UPDATE_FAILED))
}

func (u userService) UpdateProfileExtended(input *dto.ProfileEditExtendedInput) status.Object {
  profile := mapper.MapProfileUpdateExtendedInput(input)
  err := u.repo.UpdateProfileByUserId(&profile)
  return status.ConditionalRepository(err, status.UPDATED, opt.New(status.PROFILE_UPDATE_FAILED))
}

func (u userService) UpdateProfile(input *dto.ProfileEditInput) status.Object {
  profile := mapper.MapProfileUpdateInput(input)
  err := u.repo.UpdateProfileByUserId(&profile)
  return status.ConditionalRepository(err, status.UPDATED, opt.New(status.PROFILE_UPDATE_FAILED))
}

func (u userService) RegisterUser(input *dto.UserRegisterInput) (dto.UserResponse, status.Object) {
  usr, err := mapper.MapUserRegisterInput(input)
  if err != nil {
    return dto.UserResponse{}, status.InternalError()
  }

  profile := mapper.MapProfileRegisterInput(&usr, input)
  err = u.repo.CreateUser(&usr, &profile)
  if err != nil {
    return dto.UserResponse{}, status.RepositoryError(err, opt.New(status.USER_CREATION_ALREADY_EXIST))
  }

  // Send email Verification
  verifResponse, stat := u.verifService.Request(usr.Id, users.UsageVerifyEmail)
  if stat.IsError() {
    return dto.UserResponse{}, stat
  }

  m, err := mail.NewHTML("email-verification.gohtml",
    fmt.Sprintf("%s/users/email-verif/%s", u.config.ApiDNS(1), verifResponse.Token))
  if err != nil {
    return dto.UserResponse{}, status.InternalError()
  }
  m.Subject = "Email Verification"
  m.Recipients = []string{usr.Email}
  go func() {
    err := u.mailService.SendEmail(m)
    if err.IsError() {
      log.Println("erer")
    }
  }()

  return mapper.ToUserResponse(&usr), status.Created()
}

func (u userService) FindUserById(id string) (dto.UserResponse, status.Object) {
  usr, err := u.repo.FindUserById(id)
  if err != nil {
    return dto.UserResponse{}, status.RepositoryError(err, opt.New(status.USER_NOT_FOUND))
  }
  return mapper.ToUserResponse(usr), status.Success()
}

func (u userService) FindUserByEmail(email string) (dto.UserResponse, status.Object) {
  usr, err := u.repo.FindUserByEmail(email)
  if err != nil {
    return dto.UserResponse{}, status.RepositoryError(err, opt.New(status.USER_NOT_FOUND))
  }
  return mapper.ToUserResponse(usr), status.Success()
}

func (u userService) FindUserProfileById(userId string) (dto.ProfileResponse, status.Object) {
  profile, err := u.repo.FindUserProfiles(userId)
  if err != nil {
    return dto.ProfileResponse{}, status.RepositoryError(err, opt.New(status.PROFILE_NOT_FOUND))
  }
  return mapper.ToProfileResponse(profile, u.fileService), status.Success()
}

func (u userService) ChangePassword(input *dto.ChangePasswordInput) status.Object {
  // Get user
  usr, err := u.repo.FindUserById(input.UserId)
  if err != nil {
    return status.RepositoryError(err, opt.New(status.USER_NOT_FOUND))
  }

  // Match last password
  if !usr.ValidatePassword(input.LastPassword) {
    return status.Error(status.USER_CHANGE_PASSWORD_WRONG_PASSWORD)
  }

  // Set new password
  updateUser, err := mapper.MapChangePasswordInput(input)
  if err != nil {
    return status.InternalError()
  }
  err = u.repo.UpdateUser(&updateUser)
  return status.ConditionalRepository(err, status.UPDATED, opt.New(status.USER_UPDATE_FAILED))
}

func (u userService) RequestResetPassword(input *dto.ResetPasswordRequestInput) status.Object {
  userResponse, err := u.repo.FindUserByEmail(input.Email)
  if err != nil {
    return status.RepositoryError(err, opt.New(status.USER_NOT_FOUND))
  }

  verifResponse, stat := u.verifService.Request(userResponse.Id, users.UsageResetPassword)
  if stat.IsError() {
    return stat
  }

  m, err := mail.NewHTML("forgot-password.gohtml",
    fmt.Sprintf("%s/users/reset-password/%s", u.config.DNS(), verifResponse.Token))
  if err != nil {
    return status.InternalError()
  }
  m.Subject = "Reset Password"
  m.Recipients = []string{userResponse.Email}
  go u.mailService.SendEmail(m)
  return status.Success()
}

func (u userService) ResetPassword(input *dto.ResetPasswordInput) status.Object {
  // Response token
  resp, stat := u.verifService.Response(input.Token, users.UsageResetPassword)
  if stat.IsError() {
    return stat
  }

  // Logout all devices
  stat = u.authService.LogoutDevices(resp.UserId)
  if stat.IsError() {
    return stat
  }

  // Set new password
  updateUser, err := mapper.MapResetPasswordInput(input, resp.UserId)
  if err != nil {
    return status.InternalError()
  }
  err = u.repo.UpdateUser(&updateUser)
  return status.ConditionalRepository(err, status.SUCCESS, opt.New(status.USER_UPDATE_FAILED))
}

func (u userService) RequestEmailVerification(input *dto.VerifEmailRequestInput) status.Object {
  userResponse, err := u.repo.FindUserById(input.UserId)
  if err != nil {
    return status.RepositoryError(err, opt.New(status.VERIFICATION_USER_NOT_EXISTS))
  }

  verifResponse, stat := u.verifService.Request(userResponse.Id, users.UsageVerifyEmail)
  if stat.IsError() {
    return stat
  }

  // Send email
  m, err := mail.NewHTML("email-verification.gohtml",
    fmt.Sprintf("%s/users/email-verif/%s", u.config.DNS(), verifResponse.Token))
  if err != nil {
    return status.InternalError()
  }
  m.Subject = "Email Verification"
  m.Recipients = []string{userResponse.Email}
  go u.mailService.SendEmail(m)

  return status.Success()
}

func (u userService) VerifyEmail(input *dto.VerifyEmailInput) status.Object {
  response, stat := u.verifService.Response(input.Token, users.UsageVerifyEmail)
  if stat.IsError() {
    return stat
  }

  updatedUser := mapper.MapVerifyEmailInput(response.UserId)
  err := u.repo.UpdateUser(&updatedUser)
  return status.ConditionalRepository(err, status.UPDATED, opt.New(status.USER_UPDATE_FAILED))
}
