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

func (m GenreController) ListGenre(ctx *gin.Context) {
	genres, stat := m.genreService.ListGenre()
	resp.Conditional(ctx, stat, genres, nil)
}

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

func (m GenreController) UpdateGenre(ctx *gin.Context) {
	input := dto.GenreUpdateInput{}
	stat, fieldErrors := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldErrors)
		return
	}

	stat = m.genreService.UpdateGenre(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m GenreController) DeleteGenre(ctx *gin.Context) {
	genreId := ctx.Param("genre_id")
	if len(genreId) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("genre_id"))
		return
	}
	stat := m.genreService.DeleteGenre(genreId)
	resp.Conditional(ctx, stat, nil, nil)
}
