package service

import (
  "manga-explorer/internal/common/status"
  "manga-explorer/internal/domain/users/dto"
)

type IAuthentication interface {
  Authenticate(input *dto.LoginInput) (dto.LoginResponse, status.Object)
  RefreshToken(input *dto.RefreshTokenInput) (dto.RefreshTokenResponse, status.Object)
  GetCredentials(userId string) ([]dto.CredentialResponse, status.Object)
  SelfLogout(userId, accessToken string) status.Object
  Logout(userId, credId string) status.Object
  LogoutDevices(userId string) status.Object
}
