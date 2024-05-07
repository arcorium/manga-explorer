package dto

import "github.com/gin-gonic/gin"

type GenreResponse struct {
  Id   string `json:"id"`
  Name string `json:"name"`
}

type GenreCreateInput struct {
  Name string `json:"name" binding:"required"`
}

type GenreEditInput struct {
  Id   string `uri:"genre_id" binding:"required,uuid4" swaggerignore:"true"`
  Name string `json:"new_name" binding:"required"`
}

func (g *GenreEditInput) ConstructURI(ctx *gin.Context) {
  g.Id = ctx.Param("genre_id")
}
