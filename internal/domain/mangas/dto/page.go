package dto

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type PageResponse struct {
	Page     uint16 `json:"page"`
	ImageURL string `json:"image_url"` // Can be returning image bytes
}

type InternalPage struct {
	Number uint16                `form:"page" binding:"required,gte=1"`
	Image  *multipart.FileHeader `form:"image" binding:"required"`
}

type PageCreateInput struct {
	ChapterId string `uri:"chapter_id" binding:"required,uuid4" swaggerignore:"true"`
	Page      InternalPage
	//InternalPage
}

func (p *PageCreateInput) ConstructURI(ctx *gin.Context) {
	p.ChapterId = ctx.Param("chapter_id")
}

type PageDeleteInput struct {
	ChapterId string   `uri:"chapter_id" binding:"required,uuid4" swaggerignore:"true"`
	Pages     []uint16 `json:"pages" binding:"required"`
}

func (d *PageDeleteInput) ConstructURI(ctx *gin.Context) {
	d.ChapterId = ctx.Param("chapter_id")
}
