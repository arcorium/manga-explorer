package dto

import "manga-explorer/internal/domain/users/dto"

type RateResponse struct {
	User dto.UserResponse `json:"user"`
	Rate uint8            `json:"rate"`
}

type RateUpsertInput struct {
	UserId  string `json:"-"`
	MangaId string `uri:"manga_id" binding:"required,uuid4"`
	Rate    uint8  `json:"rate" binding:"required,lte=10"`
}
