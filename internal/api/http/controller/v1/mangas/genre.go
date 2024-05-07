package mangas

import (
  "github.com/gin-gonic/gin"
  "manga-explorer/internal/common"
  "manga-explorer/internal/common/status"
  "manga-explorer/internal/domain/mangas/dto"
  "manga-explorer/internal/domain/mangas/service"
  "manga-explorer/internal/util/httputil"
  "manga-explorer/internal/util/httputil/resp"
)

func NewGenreController(genreService service.IGenre) GenreController {
  return GenreController{genreService: genreService}
}

type GenreController struct {
  genreService service.IGenre
}

// @Summary		Get All Genres
// @Description	Get all registered genres
// @Tags			manga, genre
// @Produce		json
// @Success		200	{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.GenreResponse}}
// @Router			/genres [get]
func (m GenreController) ListGenre(ctx *gin.Context) {
  genres, stat := m.genreService.ListGenre()
  resp.Conditional(ctx, stat, genres, nil)
}

// @Summary		Create Genre
// @Description	Create new genre
// @Tags			manga, genre
// @Accept			json
// @Produce		json
// @Param			input	body		dto.GenreCreateInput	true	"genre create input"
// @Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/genres [post]
func (m GenreController) CreateGenre(ctx *gin.Context) {
  genreInput := dto.GenreCreateInput{}
  stat, fieldsErr := httputil.BindJson(ctx, &genreInput)
  if stat.IsError() {
    resp.ErrorDetailed(ctx, stat, fieldsErr)
    return
  }

  stat = m.genreService.CreateGenre(genreInput)
  resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Edit Genre
// @Description	Edit specific genre by id
// @Tags			manga, genre
// @Accept			json
// @Produce		json
// @Param			input		body		dto.GenreEditInput	true	"genre edit input"
// @Param			genre_id	path		uuid.UUID			true	"genre id"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/genres/{genre_id} [put]
func (m GenreController) EditGenre(ctx *gin.Context) {
  input := dto.GenreEditInput{}
  input.ConstructURI(ctx)
  stat, fieldErrors := httputil.BindJson(ctx, &input)
  if stat.IsError() {
    resp.ErrorDetailed(ctx, stat, fieldErrors)
    return
  }

  stat = m.genreService.UpdateGenre(&input)
  resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Delete Genre
// @Description	Delete specific genre by id
// @Tags			manga, genre
// @Produce		json
// @Param			genre_id	path		uuid.UUID	true	"genre id"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Router			/genres/{genre_id} [delete]
func (m GenreController) DeleteGenre(ctx *gin.Context) {
  genreId := ctx.Param("genre_id")
  if len(genreId) == 0 {
    resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("genre_id"))
    return
  }
  stat := m.genreService.DeleteGenre(genreId)
  resp.Conditional(ctx, stat, nil, nil)
}
