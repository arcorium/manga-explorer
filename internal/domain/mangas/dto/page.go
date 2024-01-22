package dto

import "mime/multipart"

type PageResponse struct {
	Page     uint16 `json:"page"`
	ImageURL string `json:"image_url"` // Can be returning image bytes
}

type PageCreateInput struct {
	ChapterId string                `uri:"chapter_id" binding:"required,uuid4"`
	Page      uint16                `form:"page" binding:"required"`
	Image     *multipart.FileHeader `form:"image" binding:"required"`
}

type PageDeleteInput struct {
	ChapterId string   `uri:"chapter_id" binding:"required,uuid4"`
	Pages     []uint16 `json:"pages" binding:"required"`
}
