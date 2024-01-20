package mangas

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/app/dto"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/httputil"
)

func NewMangaController(mangaService service.IManga) MangaController {
	return MangaController{mangaService: mangaService}
}

type MangaController struct {
	mangaService service.IManga
}

func (m MangaController) ListManga(ctx *gin.Context) {
	var pagedQuery dto.PagedQueryInput
	if err := ctx.BindQuery(&pagedQuery); err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_QUERY_ERROR))
		return
	}

	if pagedQuery.IsPaged() {
		mangas, pages, status := m.mangaService.ListPagedMangas(&pagedQuery)
		httputil.PagedResponse(ctx, status, &pages, mangas)
	} else {
		mangas, status := m.mangaService.ListMangas()
		httputil.Response(ctx, status, mangas)
	}
}

func (m MangaController) Search(ctx *gin.Context) {
	searchQuery, status := httputil.BindQueryJson[mangaDto.MangaSearchQuery](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	if searchQuery.IsPaged() {
		mangas, page, status := m.mangaService.SearchPagedMangas(&searchQuery)
		httputil.PagedResponse(ctx, status, &page, mangas)
	} else {
		mangas, status := m.mangaService.SearchMangas(&searchQuery)
		httputil.Response(ctx, status, mangas)
	}
}

func (m MangaController) EditManga(ctx *gin.Context) {
	updateManga, status := httputil.BindUriJson[mangaDto.MangaEditInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = m.mangaService.EditManga(&updateManga)
	httputil.Response(ctx, status, nil)
}

func (m MangaController) Random(ctx *gin.Context) {
	limit := util.GetDefaultedUintQuery(ctx, "limit", 1)
	mangas, status := m.mangaService.FindRandomMangas(limit)
	httputil.Response(ctx, status, mangas)
}

func (m MangaController) FindMangaById(ctx *gin.Context) {
	id := ctx.Param("manga_id")
	if len(id) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}

	mangas, status := m.mangaService.FindMangaByIds(id)
	httputil.Response(ctx, status, mangas)
}

func (m MangaController) FindMangaComments(ctx *gin.Context) {
	id := ctx.Param("manga_id")
	if len(id) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}

	comments, status := m.mangaService.FindMangaComments(id)
	httputil.Response(ctx, status, comments)
}

func (m MangaController) FindMangaRatings(ctx *gin.Context) {
	id := ctx.Param("manga_id")
	if len(id) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}

	rates, status := m.mangaService.FindMangaRatings(id)
	httputil.Response(ctx, status, rates)
}

func (m MangaController) CreateMangaComments(ctx *gin.Context) {
	commentInput, status := httputil.BindUriJson[mangaDto.MangaCommentCreateInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = m.mangaService.CreateComments(&commentInput)
	httputil.Response(ctx, status, nil)
}

func (m MangaController) CreateMangaRatings(ctx *gin.Context) {
	rateInput, status := httputil.BindUriJson[mangaDto.RateUpsertInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = m.mangaService.UpsertMangaRating(&rateInput)
	httputil.Response(ctx, status, nil)
}

func (m MangaController) CreateManga(ctx *gin.Context) {
	var mangaInput mangaDto.MangaCreateInput
	if err := ctx.BindJSON(&mangaInput); err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_BODY_REQUEST_ERROR))
		return
	}

	status := m.mangaService.CreateManga(&mangaInput)
	httputil.Response(ctx, status, nil)
}

func (m MangaController) CreateVolume(ctx *gin.Context) {
	volumeInput, status := httputil.BindUriJson[mangaDto.VolumeCreateInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = m.mangaService.CreateVolume(&volumeInput)
	httputil.Response(ctx, status, nil)
}

func (m MangaController) DeleteVolume(ctx *gin.Context) {
	var volumeInput mangaDto.VolumeDeleteInput
	if err := ctx.BindUri(volumeInput); err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}

	status := m.mangaService.DeleteVolume(&volumeInput)
	httputil.Response(ctx, status, nil)
}

func (m MangaController) GetMangaHistories(ctx *gin.Context) {
	var pagedQuery dto.PagedQueryInput
	if err := ctx.BindQuery(&pagedQuery); err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_QUERY_ERROR))
		return
	}

	claims, err := util.GetContextValue[*common.AccessTokenClaims](ctx, common.ClaimsKey)
	if err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.AUTH_UNAUTHORIZED))
		return
	}

	if pagedQuery.IsPaged() {
		mangas, pages, cerr := m.mangaService.FindPagedMangaHistories(claims.UserId, &pagedQuery)
		httputil.PagedResponse(ctx, cerr, &pages, mangas)
	} else {
		mangas, cerr := m.mangaService.FindMangaHistories(claims.UserId)
		httputil.Response(ctx, cerr, mangas)
	}
}

func (m MangaController) GetMangaFavorites(ctx *gin.Context) {
	var pagedQuery dto.PagedQueryInput
	if err := ctx.BindQuery(&pagedQuery); err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_QUERY_ERROR))
		return
	}

	claims, err := util.GetContextValue[*common.AccessTokenClaims](ctx, common.ClaimsKey)
	if err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.AUTH_UNAUTHORIZED))
		return
	}

	if pagedQuery.IsPaged() {
		mangas, pages, cerr := m.mangaService.FindPagedMangaFavorites(claims.UserId, &pagedQuery)
		httputil.PagedResponse(ctx, cerr, &pages, mangas)
	} else {
		mangas, cerr := m.mangaService.FindMangaFavorites(claims.UserId)
		httputil.Response(ctx, cerr, mangas)
	}
}
