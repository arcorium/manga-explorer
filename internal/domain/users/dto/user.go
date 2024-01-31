package dto

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/common/status"
)

type UserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserRegisterInput struct {
	Username  string `json:"username" binding:"required,gte=5"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,gte=7"`
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname"`
	Bio       string `json:"bio"`
}

type ResetPasswordRequestInput struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifEmailRequestInput struct {
	UserId string `json:"-"`
}

type VerifyEmailInput struct {
	Token string `uri:"token" binding:"required"`
}

func (v *VerifyEmailInput) ConstructURI(ctx *gin.Context) {
	v.Token = ctx.Param("token")
}

type UserUpdateInput struct {
	UserId   string `json:"-"`
	Username string `json:"username" binding:"required,gt=5"`
	Email    string `json:"email" binding:"required,email"`
}

func (u *UserUpdateInput) Status() status.Object {
	if len(u.Username) == 0 && len(u.Email) == 0 {
		return status.ErrorMessage("There is should be at least one field to be updated")
	}
	return status.InternalSuccess()
}

type ChangePasswordInput struct {
	UserId       string `json:"-"`
	LastPassword string `json:"last_password" binding:"required"`
	NewPassword  string `json:"new_password" binding:"required,gte=7"`
}

type ResetPasswordInput struct {
	Token       string `uri:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,gte=7"`
}

func (r *ResetPasswordInput) ConstructURI(ctx *gin.Context) {
	r.Token = ctx.Param("token")
}

type UserUpdateExtendedInput struct {
	UserId   string `uri:"id" binding:"required"`
	Username string `json:"username"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password"`
}

func (u *UserUpdateExtendedInput) ConstructURI(ctx *gin.Context) {
	u.UserId = ctx.Param("id")
}

type ProfileUpdateExtendedInput struct {
	UserId    string `uri:"id" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
	PhotoURL  string `json:"photo_url"`
}

func (p *ProfileUpdateExtendedInput) ConstructURI(ctx *gin.Context) {
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
