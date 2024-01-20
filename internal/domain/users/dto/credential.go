package dto

import (
	"manga-explorer/internal/app/common/constant"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9|._]+@[a-zA-Z0-9]+(\\.[a-zA-Z0-9]{2,})+$")

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenInput struct {
	Type        string `json:"token_type" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
}

var NoLoginResponse = LoginResponse{}

type LoginResponse struct {
	Type        string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

var NoRefreshTokenResponse = RefreshTokenResponse{}

type RefreshTokenResponse struct {
	Type        string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type CredentialResponse struct {
	Id         string `json:"id"`
	DeviceName string `json:"device_name"`
}

func NewLoginResponse(token string) LoginResponse {
	return LoginResponse{
		Type:        constant.TokenType,
		AccessToken: token,
	}
}

func NewRefreshTokenResponse(token string) RefreshTokenResponse {
	return RefreshTokenResponse{
		Type:        constant.TokenType,
		AccessToken: token,
	}
}
