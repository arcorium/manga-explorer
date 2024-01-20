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
	"manga-explorer/internal/util/httputil/resp"
)

func NewMangaController(mangaService service.IManga) MangaController {
	return MangaController{mangaService: mangaService}
}

type MangaController struct {
	mangaService service.IManga
}

func (m MangaController) ListManga(ctx *gin.Context) {
	var pagedQuery dto.PagedQueryInput
	stat, fieldsErr := httputil.BindQuery(ctx, &pagedQuery)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	mangas, pages, stat := m.mangaService.ListMangas(&pagedQuery)
	resp.Conditional(ctx, stat, mangas, pages)
}

func (m MangaController) Search(ctx *gin.Context) {
	var searchQuery mangaDto.MangaSearchQuery
	stat, fieldsErr := httputil.BindQueryJson(ctx, &searchQuery)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	mangas, page, stat := m.mangaService.SearchPagedMangas(&searchQuery)
	resp.Conditional(ctx, stat, mangas, page)
}

func (m MangaController) EditManga(ctx *gin.Context) {
	var updateManga mangaDto.MangaEditInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &updateManga)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.EditManga(&updateManga)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) Random(ctx *gin.Context) {
	limit := util.GetDefaultedUintQuery(ctx, "limit", 1)
	mangas, stat := m.mangaService.FindRandomMangas(limit)
	resp.Conditional(ctx, stat, mangas, nil)
}

func (m MangaController) FindMangaById(ctx *gin.Context) {
	id := ctx.Param("manga_id")
	if len(id) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("manga_id"))
		return
	}

	mangas, stat := m.mangaService.FindMangaByIds(id)
	resp.Conditional(ctx, stat, mangas, nil)
}

func (m MangaController) FindMangaComments(ctx *gin.Context) {
	id := ctx.Param("manga_id")
	if len(id) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("manga_id"))
		return
	}

	comments, stat := m.mangaService.FindMangaComments(id)
	resp.Conditional(ctx, stat, comments, nil)
}

func (m MangaController) FindMangaRatings(ctx *gin.Context) {
	id := ctx.Param("manga_id")
	if len(id) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("manga_id"))
		return
	}

	rates, stat := m.mangaService.FindMangaRatings(id)
	resp.Conditional(ctx, stat, rates, nil)
}

func (m MangaController) CreateMangaComments(ctx *gin.Context) {
	var commentInput mangaDto.MangaCommentCreateInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &commentInput)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.CreateComments(&commentInput)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) CreateMangaRatings(ctx *gin.Context) {
	var rateInput mangaDto.RateUpsertInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &rateInput)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.UpsertMangaRating(&rateInput)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) CreateManga(ctx *gin.Context) {
	var mangaInput mangaDto.MangaCreateInput
	stat, fieldsErr := httputil.BindJson(ctx, &mangaInput)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.CreateManga(&mangaInput)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) CreateVolume(ctx *gin.Context) {
	var volumeInput mangaDto.VolumeCreateInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &volumeInput)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.CreateVolume(&volumeInput)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) DeleteVolume(ctx *gin.Context) {
	var volumeInput mangaDto.VolumeDeleteInput
	stat, fieldsErr := httputil.BindUri(ctx, &volumeInput)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.DeleteVolume(&volumeInput)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) GetMangaHistories(ctx *gin.Context) {
	var pagedQuery dto.PagedQueryInput
	stat, fieldsErr := httputil.BindQuery(ctx, &pagedQuery)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}

	mangas, pages, cerr := m.mangaService.FindMangaHistories(claims.UserId, &pagedQuery)
	resp.Conditional(ctx, cerr, mangas, pages)
}

func (m MangaController) GetMangaFavorites(ctx *gin.Context) {
	var pagedQuery dto.PagedQueryInput
	stat, fieldsErr := httputil.BindQuery(ctx, &pagedQuery)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}

	mangas, pages, cerr := m.mangaService.FindMangaFavorites(claims.UserId, &pagedQuery)
	resp.Conditional(ctx, cerr, mangas, pages)
}
