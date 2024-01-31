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
	UserId  string `json:"-"`
	MangaId string `uri:"manga_id" binding:"required,uuid4"`
	Rate    uint8  `json:"rate" binding:"required,lte=10"`
}

func (r *RateUpsertInput) ConstructURI(ctx *gin.Context) {
	r.MangaId = ctx.Param("manga_id")
}
