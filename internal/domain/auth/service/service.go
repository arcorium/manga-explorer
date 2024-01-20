package service

import (
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/domain/auth/dto"
)

type IAuthentication interface {
	Authenticate(request *dto.LoginInput) (dto.LoginResponse, common.Status)
	RefreshToken(request *dto.RefreshTokenInput) (dto.RefreshTokenResponse, common.Status)
	GetCredentials(userId string) ([]dto.CredentialResponse, common.Status)
	SelfLogout(userId, accessToken string) common.Status
	Logout(userId, credId string) common.Status
	LogoutDevices(userId string) common.Status
}
