package dto

import (
	"manga-explorer/internal/app/common"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9|._]+@[a-zA-Z0-9]+(\\.[a-zA-Z0-9]{2,})+$")

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (l LoginInput) Validate() bool {
	return emailRegex.MatchString(l.Email) && len(l.Password) > 0
}

type ResetPasswordRequestInput struct {
	Email string `json:"email" binding:"required"`
}

func (r ResetPasswordRequestInput) Validate() bool {
	return emailRegex.MatchString(r.Email)
}

type RefreshTokenInput struct {
	Type        string `json:"token_type" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
}

func (r RefreshTokenInput) Validate() bool {
	return r.Type == common.TokenType && len(r.AccessToken) > 0
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
		Type:        common.TokenType,
		AccessToken: token,
	}
}

func NewRefreshTokenResponse(token string) RefreshTokenResponse {
	return RefreshTokenResponse{
		Type:        common.TokenType,
		AccessToken: token,
	}
}
