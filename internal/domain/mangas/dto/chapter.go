package dto

import (
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/domain/users/dto"
	"time"
)

type ChapterResponse struct {
	Language  common.Country `json:"language" `
	Title     string         `json:"title"`
	CreatedAt time.Time      `json:"created_at"`

	Comments   []CommentResponse `json:"comments"`
	Pages      []PageResponse    `json:"pages"`
	Translator dto.UserResponse  `json:"translator"`
}

type ChapterCreateInput struct {
	MangaId      string         `uri:"manga_id" binding:"required,uuid4"`
	VolumeId     string         `json:"volume_id" binding:"required,uuid4"`
	Language     common.Country `json:"language" binding:"required,iso3166_1_alpha3|iso3166_1_alpha2"`
	Title        string         `json:"title" binding:"required"`
	PublishDate  time.Time      `json:"publish_date"`
	TranslatorId string         `json:"-"`
}

type ChapterEditInput struct {
	ChapterId   string         `uri:"chapter_id" binding:"required,uuid4"`
	VolumeId    string         `json:"volume_id" binding:"required,uuid4"`
	Title       string         `json:"title"`
	Language    common.Country `json:"language" binding:"iso3166_1_alpha3|iso3166_1_alpha2"`
	PublishDate time.Time      `json:"publish_date"`
}
