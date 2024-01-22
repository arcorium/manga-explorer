package dto

import (
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/dto"
	"time"
)

type MangaResponse struct {
	Id              string                `json:"id"`
	Status          uint8                 `json:"status"`
	Origin          common.Country        `json:"origin"`
	PublicationYear uint16                `json:"year"`
	Comments        []CommentResponse     `json:"comments"`
	Ratings         []RateResponse        `json:"ratings"`
	Translations    []TranslationResponse `json:"translations"`
	Volumes         []VolumeResponse      `json:"volumes"`
	ViewedCount     uint64                `json:"viewed_count"`
	FavoriteCount   uint64                `json:"favorite_count"`
}

type MangaHistoryResponse struct {
	MangaResponse
	LastView time.Time `json:"last_view"`
}

type MangaFavoriteResponse struct {
	MangaResponse
	FavoritedAt time.Time `json:"favorited_at"`
}

type MangaCreateInput struct {
	Descriptions    []TranslationCreateInput `json:"descriptions" binding:"required"`
	Status          string                   `json:"status" binding:"required,manga_status"`
	Origin          common.Country           `json:"origin" binding:"required,iso3166_1_alpha3|iso3166_1_alpha2"`
	PublicationYear uint16                   `json:"publication_year" binding:"required"`
}

type MangaEditInput struct {
	MangaId         string         `uri:"manga_id" binding:"required"`
	Status          string         `json:"status" binding:"required,manga_status"`
	Origin          common.Country `json:"origin" binding:"iso3166_1_alpha3|iso3166_1_alpha2"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	PublicationYear uint16         `json:"publication_year"`
	CoverUrl        string         `json:"cover_url"`
}

type MangaSearchQuery struct {
	dto.PagedQueryInput
	Title  string                              `json:"title"`
	Genres common.CriterionOption[string]      `json:"genre"`
	Origin common.IncludeArray[common.Country] `json:"origin"`
}
