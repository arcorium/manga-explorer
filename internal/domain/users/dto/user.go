package dto

import (
	"manga-explorer/internal/app/common"
)

type UserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserRegisterInput struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
}

type ResetPasswordRequestInput struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdateUserInput struct {
	UserId   string `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password"`
}

func (u *UpdateUserInput) SetUserId(claims *common.AccessTokenClaims) {
	u.UserId = claims.UserId
}

type ChangePasswordInput struct {
	UserId       string `json:"-"`
	LastPassword string `json:"last_password"`
	NewPassword  string `json:"new_password"`
}

func (c *ChangePasswordInput) SetUserId(claims *common.AccessTokenClaims) {
	c.UserId = claims.UserId
}

type ResetPasswordInput struct {
	Token       string `uri:"token"`
	UserId      string `json:"-"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdateUserExtendedInput struct {
	UserId   string `uri:"id" binding:"required"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileExtendedInput struct {
	UserId    string `uri:"id" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
	PhotoURL  string `json:"photo_url"`
}

type AddUserInput struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Verified  bool   `json:"verified"`
	Role      uint8  `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
	PhotoURL  string `json:"photo_url"`
}
