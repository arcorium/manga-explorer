package dto

import "mime/multipart"

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
}

type ProfileImageUpdateInput struct {
	UserId string                `json:"-"`
	Image  *multipart.FileHeader `form:"image"` // When null, it is considered to be delete
}
