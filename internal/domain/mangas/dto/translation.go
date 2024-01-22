package dto

import "manga-explorer/internal/app/common"

type TranslationCreateInput struct {
	Lang        common.Country `json:"lang" binding:"iso3166_1_alpha3|iso3166_1_alpha2"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
}

type TranslationResponse struct {
	Lang        common.Language `json:"lang"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
}
