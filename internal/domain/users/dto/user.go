package dto

import (
	"github.com/gin-gonic/gin"
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

type ChangePasswordInput struct {
	UserId       string `json:"-"`
	LastPassword string `json:"last_password"`
	NewPassword  string `json:"new_password"`
}

type ResetPasswordInput struct {
	Token       string `uri:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (r *ResetPasswordInput) ConstructURI(ctx *gin.Context) {
	r.Token = ctx.Param("token")
}

type UpdateUserExtendedInput struct {
	UserId   string `uri:"id" binding:"required"`
	Username string `json:"username"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password"`
}

func (u *UpdateUserExtendedInput) ConstructURI(ctx *gin.Context) {
	u.UserId = ctx.Param("id")
}

type UpdateProfileExtendedInput struct {
	UserId    string `uri:"id" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
	PhotoURL  string `json:"photo_url"`
}

func (p *UpdateProfileExtendedInput) ConstructURI(ctx *gin.Context) {
	p.UserId = ctx.Param("id")
}

type AddUserInput struct {
	Username  string `json:"username"`
	Email     string `json:"email" binding:"email"`
	Password  string `json:"password"`
	Verified  bool   `json:"verified"`
	Role      string `json:"role" binding:"oneof= admin user"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
	PhotoURL  string `json:"photo_url"`
}
