package dto

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/common"
	"manga-explorer/internal/domain/users/dto"
	"time"
)

type ChapterResponse struct {
	Id           string          `json:"id"`
	Language     common.Language `json:"language"`
	Chapter      uint64          `json:"chapter"`
	Title        string          `json:"title"`
	CreatedAt    time.Time       `json:"created_at"`
	TotalComment *uint64         `json:"total_comment,omitempty"`

	Comments   []CommentResponse `json:"comments,omitempty"`
	Pages      []PageResponse    `json:"pages,omitempty"`
	Translator dto.UserResponse  `json:"translator,omitempty"`
}

type ChapterCreateInput struct {
	MangaId      string          `uri:"manga_id" binding:"required,uuid4"`
	VolumeId     string          `json:"volume_id" binding:"required,uuid4"`
	Language     common.Language `json:"language" binding:"required,language"`
	Title        string          `json:"title" binding:"required"`
	Chapter      uint64          `json:"chapter" binding:"required,min=1"`
	PublishDate  time.Time       `json:"publish_date"`
	TranslatorId string          `json:"-"`
}

func (c *ChapterCreateInput) ConstructURI(ctx *gin.Context) {
	c.MangaId = ctx.Param("manga_id")
}

type ChapterEditInput struct {
	ChapterId   string          `uri:"chapter_id" binding:"required,uuid4"`
	VolumeId    string          `json:"volume_id" binding:"required,uuid4"`
	Title       string          `json:"title" binding:"required"`
	Language    common.Language `json:"language" binding:"required,language"`
	Chapter     uint64          `json:"chapter" binding:"required"`
	PublishDate time.Time       `json:"publish_date"`
}

func (e *ChapterEditInput) ConstructURI(ctx *gin.Context) {
	res := ctx.Param("chapter_id")
	e.ChapterId = res
}
