package dto

import (
  "github.com/gin-gonic/gin"
  "manga-explorer/internal/domain/users/dto"
  "time"
)

type RateResponse struct {
  User dto.UserResponse `json:"user"`
  Rate uint8            `json:"rate"`
  Time time.Time        `json:"time"`
}

type RateUpsertInput struct {
  UserId  string `json:"-" swaggerignore:"true"`
  MangaId string `uri:"manga_id" binding:"required,uuid4" swaggerignore:"true"`
  Rate    uint8  `json:"rate" binding:"required,lte=10"`
}

func (r *RateUpsertInput) ConstructURI(ctx *gin.Context) {
  r.MangaId = ctx.Param("manga_id")
}
