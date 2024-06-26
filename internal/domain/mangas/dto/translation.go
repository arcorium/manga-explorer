package dto

import (
  "github.com/gin-gonic/gin"
  "manga-explorer/internal/common"
)

type InternalTranslation struct {
  Lang        common.Language `json:"lang" binding:"required,bcp47_language_tag"`
  Title       string          `json:"title" binding:"required"`
  Description string          `json:"desc" binding:"required"`
}

type MangaTranslationInsertInput struct {
  MangaId      string                `uri:"manga_id" binding:"required,uuid4" swaggerignore:"true"`
  Translations []InternalTranslation `json:"translations" binding:"min=1"`
}

func (i *MangaTranslationInsertInput) ConstructURI(ctx *gin.Context) {
  i.MangaId = ctx.Param("manga_id")
}

type TranslationDeleteInput struct {
  TranslationIds []string `json:"ids" binding:"required,dive,uuid4"`
}

type TranslationEditInput struct {
  TranslationId string          `uri:"translate_id" binding:"required,uuid4" swaggerignore:"true"`
  Lang          common.Language `json:"lang" binding:"required,bcp47_language_tag"`
  Title         string          `json:"title" binding:"required"`
  Description   string          `json:"desc" binding:"required"`
}

func (t *TranslationEditInput) ConstructURI(ctx *gin.Context) {
  t.TranslationId = ctx.Param("translate_id")
}

type MangaTranslationsDeleteInput struct {
  MangaId   string            `uri:"manga_id" binding:"required,uuid4" swaggerignore:"true"`
  Languages []common.Language `json:"lang" binding:"required,dive,bcp47_language_tag"`
}

func (t *MangaTranslationsDeleteInput) ConstructURI(ctx *gin.Context) {
  t.MangaId = ctx.Param("manga_id")
}

type TranslationResponse struct {
  Id          string          `json:"id"`
  Lang        common.Language `json:"lang"`
  Title       string          `json:"title"`
  Description string          `json:"desc"`
}
