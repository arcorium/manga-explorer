package dto

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"strconv"
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
	Pages     []InternalPage
	//InternalPage
}

func (p *PageCreateInput) ConstructURI(ctx *gin.Context) {
	p.ChapterId = ctx.Param("chapter_id")
}

func (p *PageCreateInput) ParseMultipart(ctx *gin.Context) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	for k, v := range form.File {
		page, err := strconv.ParseUint(k, 10, 16)
		if err != nil {
			return err
		}

		// check if one page contains multiple images
		if len(v) > 1 {
			return errors.New("each page should only have one image")
		}
		p.Pages = append(p.Pages, InternalPage{Number: uint16(page), Image: v[0]})
	}

	return nil
}

type PageDeleteInput struct {
	ChapterId string   `uri:"chapter_id" binding:"required,uuid4" swaggerignore:"true"`
	Pages     []uint16 `json:"pages" binding:"required"`
}

func (d *PageDeleteInput) ConstructURI(ctx *gin.Context) {
	d.ChapterId = ctx.Param("chapter_id")
}
