package dto

import "github.com/gin-gonic/gin"

type VolumeResponse struct {
	Title    string            `json:"title"`
	Number   uint32            `json:"number"`
	Chapters []ChapterResponse `json:"chapters"`
}

type VolumeCreateInput struct {
	MangaId string `uri:"manga_id" binding:"required,uuid4"`
	Title   string `json:"title"`
	Number  uint32 `json:"number" binding:"required"`
}

func (c *VolumeCreateInput) ConstructURI(ctx *gin.Context) {
	c.MangaId = ctx.Param("manga_id")
}

type VolumeDeleteInput struct {
	MangaId string `uri:"manga_id" binding:"required,uuid4"`
	Volume  uint32 `uri:"volume" binding:"required"`
}
