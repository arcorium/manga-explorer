package mangas

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/util/httputil"
)

func NewGenreController(genreService service.IGenre) GenreController {
	return GenreController{genreService: genreService}
}

type GenreController struct {
	genreService service.IGenre
}

func (m GenreController) ListGenre(ctx *gin.Context) {
	genres, status := m.genreService.ListGenre()
	httputil.Response(ctx, status, genres)
}

func (m GenreController) CreateGenre(ctx *gin.Context) {
	var genreInput dto.GenreCreateInput
	if err := ctx.BindJSON(&genreInput); err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_BODY_REQUEST_ERROR))
		return
	}

	status := m.genreService.CreateGenre(genreInput)
	httputil.Response(ctx, status, nil)
}

func (m GenreController) DeleteGenre(ctx *gin.Context) {
	genreId := ctx.Param("genre_id")
	if len(genreId) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}
	status := m.genreService.DeleteGenre(genreId)
	httputil.Response(ctx, status, nil)
}
