package dto

import (
	"manga-explorer/internal/app/common"
)

type ProfileResponse struct {
	UserResponse UserResponse `json:"user"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PhotoURL  string `json:"photo_url"`
	Bio       string `json:"bio"`
}

type ProfileUpdateInput struct {
	UserId    string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
	PhotoURL  string `json:"photo_url"`
}

func (p *ProfileUpdateInput) SetUserId(claims *common.AccessTokenClaims) {
	p.UserId = claims.UserId
}
