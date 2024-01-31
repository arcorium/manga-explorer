package dto

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/domain/users/dto"
	"time"
)

type CommentResponse struct {
	Id       string            `json:"id"`
	User     dto.UserResponse  `json:"user"`
	Comment  string            `json:"comment"`
	IsEdited bool              `json:"is_edited"`
	Like     uint64            `json:"like"`
	Dislike  uint64            `json:"dislike"`
	Time     time.Time         `json:"time"`
	Replies  []CommentResponse `json:"replies,omitempty"`
}

type commentCreatInputParent struct {
	ParentId string `json:"parent_id" binding:"omitempty,uuid4"`
	UserId   string `json:"-"`
	Comment  string `json:"comment" binding:"required"`
}

func (c commentCreatInputParent) HasParent() bool {
	return len(c.ParentId) != 0
}

type MangaCommentCreateInput struct {
	MangaId string `uri:"manga_id" binding:"required"`
	commentCreatInputParent
}

func (c *MangaCommentCreateInput) ConstructURI(ctx *gin.Context) {
	c.MangaId = ctx.Param("manga_id")
}

type ChapterCommentCreateInput struct {
	ChapterId string `uri:"chapter_id" binding:"required"`
	commentCreatInputParent
}

func (c *ChapterCommentCreateInput) ConstructURI(ctx *gin.Context) {
	c.ChapterId = ctx.Param("chapter_id")
}

type PageCommentCreateInput struct {
	PageId string `uri:"page_id" binding:"required"`
	commentCreatInputParent
}

func (c *PageCommentCreateInput) ConstructURI(ctx *gin.Context) {
	c.PageId = ctx.Param("page_id")
}
