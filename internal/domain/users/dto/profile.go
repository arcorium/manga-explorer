package dto

import (
	"mime/multipart"
)

type InternalProfileResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PhotoURL  string `json:"photo_url"`
	Bio       string `json:"bio"`
}

type ProfileResponse struct {
	UserResponse     UserResponse            `json:"user"`
	ProfileResponses InternalProfileResponse `json:"profile"`
}

type ProfileUpdateInput struct {
	UserId    string `json:"-"`
	FirstName string `json:"first_name" binding:"required,gte=5"`
	LastName  string `json:"last_name" binding:"required"`
	Bio       string `json:"bio" binding:"required"`
}

type ProfileImageUpdateInput struct {
	UserId string                `json:"-" form:"-" binding:"required,uuid4"`
	Image  *multipart.FileHeader `form:"image" binding:"required"`
}
