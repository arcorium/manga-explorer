package mangas

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/common"
	"manga-explorer/internal/common/dto"
	"manga-explorer/internal/common/status"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/httputil"
	"manga-explorer/internal/util/httputil/resp"
	"strings"
)

func NewMangaController(mangaService service.IManga) MangaController {
	return MangaController{mangaService: mangaService}
}

type MangaController struct {
	mangaService service.IManga
}

func (m MangaController) ListManga(ctx *gin.Context) {
	input := dto.PagedQueryInput{}
	stat, fieldsErr := httputil.BindQuery(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	mangas, pages, stat := m.mangaService.ListMangas(&input)
	resp.Conditional(ctx, stat, mangas, pages)
}

func (m MangaController) Search(ctx *gin.Context) {
	input := mangaDto.MangaSearchQuery{}
	input.ConstructQuery(ctx)
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	mangas, page, stat := m.mangaService.SearchMangas(&input)
	resp.Conditional(ctx, stat, mangas, page)
}

func (m MangaController) EditManga(ctx *gin.Context) {
	input := mangaDto.MangaEditInput{}
	input.ConstructURI(ctx)
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.EditManga(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) EditMangaGenres(ctx *gin.Context) {
	input := mangaDto.MangaGenreEditInput{}
	input.ConstructURI(ctx)

	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.EditMangaGenres(&input)
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

	if !util.IsUUID(id) {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR),
			common.NewParameterError("manga_id", " should be uuid type"))
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

	if !util.IsUUID(id) {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR),
			common.NewParameterError("manga_id", " should be uuid type"))
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

	if !util.IsUUID(id) {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR),
			common.NewParameterError("manga_id", " should be uuid type"))
		return
	}

	rates, stat := m.mangaService.FindMangaRatings(id)
	resp.Conditional(ctx, stat, rates, nil)
}

func (m MangaController) CreateMangaComment(ctx *gin.Context) {
	input := mangaDto.MangaCommentCreateInput{}
	input.ConstructURI(ctx)
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	input.UserId = claims.UserId

	stat = m.mangaService.CreateComments(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) CreateMangaRating(ctx *gin.Context) {
	input := mangaDto.RateUpsertInput{}
	input.ConstructURI(ctx)
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	input.UserId = claims.UserId

	stat = m.mangaService.UpsertMangaRating(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) CreateManga(ctx *gin.Context) {
	mangaInput := mangaDto.MangaCreateInput{}
	stat, fieldsErr := httputil.BindJson(ctx, &mangaInput)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.CreateManga(&mangaInput)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) UpdateMangaCover(ctx *gin.Context) {
	input := mangaDto.MangaCoverUpdateInput{}
	input.ConstructURI(ctx)

	stat, fieldsErr := httputil.BindMultipartForm(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.UpdateMangaCover(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) CreateVolume(ctx *gin.Context) {
	input := mangaDto.VolumeCreateInput{}
	input.ConstructURI(ctx)
	stat, fieldsErr := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.CreateVolume(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) InsertMangaTranslate(ctx *gin.Context) {
	input := mangaDto.MangaInsertTranslationInput{}
	input.ConstructURI(ctx)
	stat, fieldErrors := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldErrors)
		return
	}

	stat = m.mangaService.InsertMangaTranslations(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) FindMangaTranslations(ctx *gin.Context) {
	mangaId := ctx.Param("manga_id")
	if len(mangaId) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("manga_id"))
		return
	}

	if !util.IsUUID(mangaId) {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR),
			common.NewParameterError("manga_id", " should be uuid type"))
		return
	}

	language := ctx.Param("language")
	language = strings.Trim(language, "/")
	if len(language) == 0 {
		responses, stat := m.mangaService.FindMangaTranslations(mangaId)
		resp.Conditional(ctx, stat, responses, nil)
		return
	}

	lang := common.NewLanguage(language)
	if len(lang) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewParameterError("language", "language isn't supported"))
		return
	}

	translation, stat := m.mangaService.FindSpecificMangaTranslation(mangaId, lang)
	if stat.IsError() || len(translation.Id) == 0 {
		resp.Error(ctx, status.ErrorMessage("Translation not found"))
		return
	}
	resp.Success(ctx, stat, translation, nil)
}

func (m MangaController) DeleteMangaTranslations(ctx *gin.Context) {
	input := mangaDto.TranslationMangaDeleteInput{}
	input.ConstructURI(ctx)
	stat, fieldErrors := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldErrors)
		return
	}

	stat = m.mangaService.DeleteMangaTranslations(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) DeleteTranslations(ctx *gin.Context) {
	input := mangaDto.TranslationDeleteInput{}
	stat, fieldErrors := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldErrors)
		return
	}

	stat = m.mangaService.DeleteTranslations(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) UpdateTranslation(ctx *gin.Context) {
	input := mangaDto.TranslationUpdateInput{}
	input.ConstructURI(ctx)
	stat, fieldErrors := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldErrors)
		return
	}

	stat = m.mangaService.UpdateTranslation(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) DeleteVolume(ctx *gin.Context) {
	volumeInput := mangaDto.VolumeDeleteInput{}
	stat, fieldsErr := httputil.BindUri(ctx, &volumeInput)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.mangaService.DeleteVolume(&volumeInput)
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) GetMangaHistories(ctx *gin.Context) {
	query := dto.PagedQueryInput{}
	stat, fieldsErr := httputil.BindQuery(ctx, &query)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}

	mangas, pages, cerr := m.mangaService.FindMangaHistories(claims.UserId, &query)
	resp.Conditional(ctx, cerr, mangas, pages)
}

func (m MangaController) ModifyFavoriteManga(ctx *gin.Context) {
	input := mangaDto.FavoriteMangaInput{}
	stat, fieldErrors := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldErrors)
		return
	}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}
	input.UserId = claims.UserId

	if input.Operator == "add" {
		stat = m.mangaService.AddFavoriteManga(&input)
	} else if input.Operator == "remove" {
		stat = m.mangaService.RemoveFavoriteManga(&input)
	} else {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_REQUEST_ERROR), common.FieldError{Field: "op", Error: "operator should be one of add or remove"})
		return
	}
	resp.Conditional(ctx, stat, nil, nil)
}

func (m MangaController) GetMangaFavorites(ctx *gin.Context) {
	query := dto.PagedQueryInput{}
	stat, fieldsErr := httputil.BindQuery(ctx, &query)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	claims, stat := common.GetClaims(ctx)
	if stat.IsError() {
		resp.Error(ctx, stat)
		return
	}

	mangas, pages, cerr := m.mangaService.FindMangaFavorites(claims.UserId, &query)
	resp.Conditional(ctx, cerr, mangas, pages)
}
