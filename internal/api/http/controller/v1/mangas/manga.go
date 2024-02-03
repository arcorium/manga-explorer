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

// @Summary		Get All Mangas
// @Description	get all registered mangas
// @Tags			manga
// @Param			paged	query	dto.PagedQueryInput	true	"pagination query"
// @Produce		json
// @Success		200	{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.MinimalMangaResponse}}
// @Failure		400	{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas [get]
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

// @Summary		Search Manga
// @Description	search manga by the body
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			paged	query		dto.PagedQueryInput		true	"pagination query"
// @Param			input	body		dto.MangaSearchQuery	true	"search query"
// @Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.MinimalMangaResponse}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/search [get]
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

// @Summary		Edit Manga
// @Description	Edit specific manga by id
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			manga_id	path		uuid.UUID			true	"manga id"
// @Param			input		body		dto.MangaEditInput	true	"manga edit input"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id} [put]
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

// @Summary		Edit Manga Genres
// @Description	Modify specific manga's genres
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			manga_id	path		uuid.UUID				true	"manga id"
// @Param			input		body		dto.MangaGenreEditInput	true	"manga's genre edit input"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/genres [patch]
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

// @Summary		Random Manga
// @Description	Get random manga with limit query
// @Tags			manga
// @Produce		json
// @Param			limit	query		integer	false	"total response manga"
// @Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.MinimalMangaResponse}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas [get]
func (m MangaController) Random(ctx *gin.Context) {
	limit := util.GetDefaultedUintQuery(ctx, "limit", 1)
	mangas, stat := m.mangaService.FindRandomMangas(limit)
	resp.Conditional(ctx, stat, mangas, nil)
}

// @Summary	Find Manga By Id
// @Description
// @Tags		manga
// @Produce	json
// @Param		manga_id	path		uuid.UUID	true	"manga id"
// @Success	200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=dto.MangaResponse}}
// @Failure	400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
// @Failure	400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router		/mangas/{manga_id} [get]
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

// @Summary		Get Manga Comments
// @Description	Get all comments from specific manga
// @Tags			manga
// @Produce		json
// @Param			manga_id	path		uuid.UUID	true	"manga id"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.CommentResponse}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/comments [get]
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

// @Summary		Get Manga Ratings
// @Description	Get all ratings from specific manga
// @Tags			manga
// @Produce		json
// @Param			manga_id	path		uuid.UUID	true	"manga id"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.RateResponse}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/ratings [get]
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

// @Summary		Create Manga Comment
// @Description	create comment for specific manga
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			manga_id	path		uuid.UUID					true	"manga id"
// @Param			input		body		dto.MangaCommentCreateInput	true	"manga comment create input"
// @Success		201			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/comments [post]
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

// @Summary		Create Manga Rating
// @Description	create or update rate for specific manga on current logged-in user
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			manga_id	path		uuid.UUID			true	"manga id"
// @Param			input		body		dto.RateUpsertInput	true	"rate create input"
// @Success		201			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/ratings [post]
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

// @Summary		Create Manga
// @Description	create new manga
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			input	body		dto.MangaCreateInput	true	"manga create input"
// @Success		201		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas [post]
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

// @Summary		Update Manga Cover
// @Description	update manga cover image
// @Tags			manga
// @Accept			mpfd
// @Produce		json
// @Param			manga_id	path		uuid.UUID	true	"manga id"
// @Param			image		formData	file		true	"cover image"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/covers [patch]
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

// @Summary		Create Volume
// @Description	create new volume on specific manga
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			manga_id	path		uuid.UUID				true	"manga id"
// @Param			input		body		dto.VolumeCreateInput	true	"volume create input"
// @Success		201			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/volumes [post]
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

// @Summary		Insert Manga Translation
// @Description	create manga translation for specific manga
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			manga_id	path		uuid.UUID						true	"manga id"
// @Param			input		body		dto.MangaTranslationInsertInput	true	"manga translation insert input"
// @Success		201			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/translates [post]
func (m MangaController) InsertMangaTranslate(ctx *gin.Context) {
	input := mangaDto.MangaTranslationInsertInput{}
	input.ConstructURI(ctx)
	stat, fieldErrors := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldErrors)
		return
	}

	stat = m.mangaService.InsertMangaTranslations(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Find Manga Translations
// @Description	get all specifc manga translations
// @Tags			manga
// @Produce		json
// @Param			language	path		string		false	"expected translation language"
// @Param			manga_id	path		uuid.UUID	true	"manga id"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.TranslationResponse}}
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=dto.TranslationResponse}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/translates/{language} [get]
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

// @Summary		Delete Manga Translation
// @Description	delete all translations of specific manga based on languages provided
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			manga_id	path		uuid.UUID							true	"manga id"
// @Param			input		body		dto.MangaTranslationsDeleteInput	true	"manga translations delete input"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/translates [delete]
func (m MangaController) DeleteMangaTranslations(ctx *gin.Context) {
	input := mangaDto.MangaTranslationsDeleteInput{}
	input.ConstructURI(ctx)
	stat, fieldErrors := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldErrors)
		return
	}

	stat = m.mangaService.DeleteMangaTranslations(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Delete Translations
// @Description	delete specific translations based on the id
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			input	body		dto.TranslationDeleteInput	true	"translation delete input"
// @Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/translates [delete]
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

// @Summary		Edit Translation
// @Description	edit specific manga translation
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			translation_id	path		uuid.UUID					true	"translation id"
// @Param			input			body		dto.TranslationEditInput	true	"translation edit input"
// @Success		200				{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400				{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400				{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/translates/{translate_id} [put]
func (m MangaController) EditTranslation(ctx *gin.Context) {
	input := mangaDto.TranslationEditInput{}
	input.ConstructURI(ctx)
	stat, fieldErrors := httputil.BindJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldErrors)
		return
	}

	stat = m.mangaService.UpdateTranslation(&input)
	resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Delete Volume
// @Description	delete specific volumes on manga
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			manga_id	path		uuid.UUID				true	"manga id"
// @Param			input		body		dto.VolumeDeleteInput	true	"volume delete input"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/volumes [delete]
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

// @Summary		Get Manga Histories
// @Description	get current logged-in manga histories
// @Tags			manga
// @Produce		json
// @Param			page	query		dto.PagedQueryInput	false	"pagination query"
// @Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.MangaHistoryResponse}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/histories [get]
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

// @Summary		Modify Favorite Manga
// @Description	add or remove manga as favorite based on op field
// @Tags			manga
// @Accept			json
// @Produce		json
// @Param			manga_id	path		uuid.UUID							true	"manga id"
// @Param			input		body		dto.FavoriteMangaModificationInput	true	"modification input"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id}/favorites [post]
func (m MangaController) ModifyFavoriteManga(ctx *gin.Context) {
	input := mangaDto.FavoriteMangaModificationInput{}
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

// @Summary		Get Manga Favorite
// @Description	get current logged-in user manga's favorite
// @Tags			manga
// @Produce		json
// @Param			page	query		dto.PagedQueryInput	false	"pagination query"
// @Success		200		{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.MangaFavoriteResponse}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400		{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/favorites [get]
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
