package dto

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
)

type InternalTranslation struct {
	Lang        common.Country `json:"lang" binding:"iso3166_1_alpha3|iso3166_1_alpha2"`
	Title       string         `json:"title"`
	Description string         `json:"desc"`
}

type MangaInsertTranslationInput struct {
	MangaId      string                `uri:"manga_id" binding:"required,uuid4"`
	Translations []InternalTranslation `json:"translations"`
}

func (i *MangaInsertTranslationInput) ConstructURI(ctx *gin.Context) {
	i.MangaId = ctx.Param("manga_id")
}

type TranslationDeleteInput struct {
	TranslationIds []string `json:"ids" binding:"required,uuid4"`
}

type TranslationUpdateInput struct {
	TranslationId string         `json:"translation_id" binding:"required,uuid4"`
	Lang          common.Country `json:"lang"`
	Title         string         `json:"title"`
	Description   string         `json:"desc"`
}

type TranslationResponse struct {
	Lang        common.Language `json:"lang"`
	Title       string          `json:"title"`
	Description string          `json:"desc"`
}
