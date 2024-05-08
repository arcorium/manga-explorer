package mangas

import (
  "github.com/gin-gonic/gin"
  "manga-explorer/internal/common"
  "manga-explorer/internal/common/status"
  "manga-explorer/internal/domain/mangas/dto"
  "manga-explorer/internal/domain/mangas/service"
  "manga-explorer/internal/util"
  "manga-explorer/internal/util/httputil"
  "manga-explorer/internal/util/httputil/resp"
  "manga-explorer/internal/util/opt"
)

func NewChapterController(chapterService service.IChapter) ChapterController {
  return ChapterController{
    chapterService: chapterService,
  }
}

type ChapterController struct {
  chapterService service.IChapter
}

// @Summary		Insert Chapter Page
// @Description	insert new page for specific chapter
// @Tags			manga, chapter
// @Accept			mpfd
// @Produce		json
// @Param			chapter_id	path		uuid.UUID	true	"chapter id"
// @Param			page		formData	integer		true	"page number"
// @Param			image		formData	file		true	"page image"
// @Success		201			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/chapters/{chapter_id}/pages [post]
func (m ChapterController) InsertChapterPage(ctx *gin.Context) {
  input := dto.PageCreateInput{}
  input.ConstructURI(ctx)
  err := input.ParseMultipart(ctx)
  if err != nil {
    resp.Error(ctx, status.ErrorMessage(err.Error()))
    return
  }
  if len(input.Pages) == 0 {
    resp.Error(ctx, status.ErrorMessage("at least one page is required"))
    return
  }

  stat, errorPages := m.chapterService.InsertChapterPage(&input)
  if stat.IsError() {
    details := struct {
      Pages []uint16
    }{
      Pages: errorPages,
    }
    resp.ErrorDetailed(ctx, stat, details)
    return
  }
  resp.Success(ctx, stat, nil, nil)
}

// @Summary		Delete Chapter Pages
// @Description	delete specifc chapter pages
// @Tags			manga, chapter
// @Accept			json
// @Produce		json
// @Param			chapter_id	path		uuid.UUID			true	"chapter id"
// @Param			input		body		dto.PageDeleteInput	true	"page delete input"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/chapters/{chapter_id}/pages [delete]
func (m ChapterController) DeleteChapterPages(ctx *gin.Context) {
  pageInput := dto.PageDeleteInput{}
  pageInput.ConstructURI(ctx)
  stat, fieldsErr := httputil.BindJson(ctx, &pageInput)
  if stat.IsError() {
    resp.ErrorDetailed(ctx, stat, fieldsErr)
    return
  }

  stat = m.chapterService.DeleteChapterPages(&pageInput)
  resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Edit Chapter
// @Description	edit specific chapter
// @Tags			manga, chapter
// @Accept			json
// @Produce		json
// @Param			chapter_id	path		uuid.UUID				true	"chapter id"
// @Param			input		body		dto.ChapterEditInput	true	"chapter edit input"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/chapters/{chapter_id} [put]
func (m ChapterController) EditChapter(ctx *gin.Context) {
  editInput := dto.ChapterEditInput{}
  editInput.ConstructURI(ctx)
  stat, fieldsErr := httputil.BindJson(ctx, &editInput)
  if stat.IsError() {
    resp.ErrorDetailed(ctx, stat, fieldsErr)
    return
  }

  stat = m.chapterService.EditChapter(&editInput)
  resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Create Chapter Comment
// @Description	create comment for specific chapter
// @Tags			manga, chapter
// @Produce		json
// @Param			chapter_id	path		uuid.UUID	true	"chapter id"
// @Success		201			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/chapters/{chapter_id}/comments [post]
func (m ChapterController) CreateChapterComments(ctx *gin.Context) {
  input := dto.ChapterCommentCreateInput{}
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

  stat = m.chapterService.CreateChapterComment(&input)
  resp.Conditional(ctx, stat, nil, nil)
}

//func (m ChapterController) CreatePageComments(ctx *gin.Context) {
//	input := dto.PageCommentCreateInput{}
//	input.ConstructURI(ctx)
//	stat, fieldsErr := httputil.BindJson(ctx, &input)
//	if stat.IsError() {
//		resp.ErrorDetailed(ctx, stat, fieldsErr)
//		return
//	}
//	stat = m.chapterService.CreatePageComment(&input)
//	resp.Conditional(ctx, stat, nil, nil)
//}

// @Summary		Find Chapter Comment
// @Description	get all comments from specific chapter
// @Tags			manga, chapter
// @Produce		json
// @Param			chapter_id	path		uuid.UUID	true	"chapter id"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.CommentResponse}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/chapters/{chapter_id}/comments [get]
func (m ChapterController) FindChapterComments(ctx *gin.Context) {
  chapterId := ctx.Param("chapter_id")
  if len(chapterId) == 0 {
    resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("chapter_id"))
    return
  }

  if !util.IsUUID(chapterId) {
    resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.ParameterError{
      Param: "chapter_id",
      Error: "chapter_id parameter should be on uuid type",
    })
    return
  }

  comments, stat := m.chapterService.FindChapterComments(chapterId)
  resp.Conditional(ctx, stat, comments, nil)
}

//func (m ChapterController) FindPageComments(ctx *gin.Context) {
//	pageId := ctx.Param("page_id")
//	if len(pageId) == 0 {
//		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("page_id"))
//		return
//	}
//
//	if !util.IsUUID(pageId) {
//		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.ParameterError{
//			Param: "pageId",
//			Error: "page_id parameter should be on uuid type",
//		})
//		return
//	}
//
//	comments, stat := m.chapterService.FindPageComments(pageId)
//	resp.Conditional(ctx, stat, comments, nil)
//}

// @Summary		Create Chapter
// @Description	create chapter for specific manga
// @Tags			manga, chapter
// @Accept			json
// @Produce		json
// @Param			manga_id	path		uuid.UUID				true	"manga id"
// @Param			input		body		dto.ChapterCreateInput	true	"chapter create input"
// @Success		201			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/mangas/{manga_id/chapters [post]
func (m ChapterController) CreateChapter(ctx *gin.Context) {
  input := dto.ChapterCreateInput{}
  input.ConstructURI(ctx)
  stat, fieldsErr := httputil.BindJson(ctx, &input)
  if stat.IsError() {
    resp.ErrorDetailed(ctx, stat, fieldsErr)
    return
  }

  // Set translator / chapter creator
  claims, stat := common.GetClaims(ctx)
  if stat.IsError() {
    resp.Error(ctx, stat)
    return
  }
  input.TranslatorId = claims.UserId

  stat = m.chapterService.CreateChapter(&input)
  resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Delete Chapter
// @Description	delete specific manga chapter and pages associated with it
// @Tags			manga, chapter
// @Produce		json
// @Param			chapter_id	path		uuid.UUID	true	"chapter id"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=nil}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/chapters/{chapter_id} [delete]
func (m ChapterController) DeleteChapter(ctx *gin.Context) {
  chapterId := ctx.Param("chapter_id")
  if len(chapterId) == 0 {
    resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("chapter_id"))
    return
  }

  if !util.IsUUID(chapterId) {
    resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.ParameterError{
      Param: "chapter_id",
      Error: "chapter_id parameter should be on uuid type",
    })
    return
  }

  stat := m.chapterService.DeleteChapter(chapterId)
  resp.Conditional(ctx, stat, nil, nil)
}

// @Summary		Find Chapter Details
// @Description	get specific chapter details with pages associated with it
// @Tags			manga, chapter
// @Produce		json
// @Param			chapter_id	path		uuid.UUID	true	"chapter id"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=dto.ChapterResponse}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/chapters/{chapter_id} [get]
func (m ChapterController) FindChapterDetails(ctx *gin.Context) {
  chapterId := ctx.Param("chapter_id")
  if len(chapterId) == 0 {
    resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("chapter_id"))
    return
  }

  if !util.IsUUID(chapterId) {
    resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.FieldError{
      Field: "chapter_id",
      Error: "chapter_id field should be on uuid type",
    })
    return
  }

  var userId opt.Optional[string]
  claims, stat := common.GetClaims(ctx)
  if stat.IsError() {
    userId = opt.NullStr
  } else {
    userId = opt.New(claims.UserId)
  }

  pages, stat := m.chapterService.FindChapterDetails(chapterId, userId)
  resp.Conditional(ctx, stat, pages, nil)
}

// @Summary		Find Volume Details
// @Description	Get specific volume details with the chapters associated with it
// @Tags			manga, chapter
// @Produce		json
// @Param			volume_id	path		uuid.UUID	true	"volume id"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=dto.VolumeResponse}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=common.ParameterError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/volumes/{volume_id} [get]
func (m ChapterController) FindVolumeDetails(ctx *gin.Context) {
  volumeId := ctx.Param("volume_id")
  if len(volumeId) == 0 {
    resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("volume_id"))
    return
  }

  chapters, stat := m.chapterService.FindVolumeDetails(volumeId)
  resp.Conditional(ctx, stat, chapters, nil)
}

// @Summary		Get Manga History Chapters
// @Description	get chapters of manga history on current logged-in user
// @Tags			manga, chapter
// @Produce		json
// @Param			manga_id	path		uuid.UUID			true	"manga id"
// @Param			page		query		dto.PagedQueryInput	false	"pagination query"
// @Success		200			{object}	dto.SuccessWrapper{success=dto.SuccessResponse{data=[]dto.ChapterResponse}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=[]common.FieldError}}
// @Failure		400			{object}	dto.ErrorWrapper{error=dto.ErrorResponse{details=nil}}
// @Router			/chapters/{manga_id}/histories [get]
func (m ChapterController) GetMangaHistoryChapter(ctx *gin.Context) {
  input := dto.MangaChapterHistoriesFindInput{}
  input.ConstructURI(ctx)

  stat, fieldErrors := httputil.BindQuery(ctx, &input)
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

  histories, page, stat := m.chapterService.FindMangaChapterHistories(&input)
  resp.Conditional(ctx, stat, histories, page)
}
