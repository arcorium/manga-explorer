package mangas

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	dto2 "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/httputil"
)

func NewChapterController(chapterService service.IChapter) ChapterController {
	return ChapterController{
		chapterService: chapterService,
	}
}

type ChapterController struct {
	chapterService service.IChapter
}

func (m ChapterController) InsertChapterPage(ctx *gin.Context) {
	pageInput, status := httputil.BindUriMultipartForm[dto2.PageCreateInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}
	status = m.chapterService.InsertChapterPage(&pageInput)
	httputil.Response(ctx, status, nil)
}
func (m ChapterController) DeleteChapterPage(ctx *gin.Context) {
	pageInput, status := httputil.BindUriJson[dto2.PageDeleteInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
	}

	status = m.chapterService.DeleteChapterPages(&pageInput)
	httputil.Response(ctx, status, nil)
}
func (m ChapterController) EditChapter(ctx *gin.Context) {
	editInput, status := httputil.BindUriJson[dto2.ChapterEditInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = m.chapterService.EditChapter(&editInput)
	httputil.Response(ctx, status, nil)
}
func (m ChapterController) CreateChapterComments(ctx *gin.Context) {
	commentInput, status := httputil.BindUriJson[dto2.ChapterCommentCreateInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}

	status = m.chapterService.CreateChapterComment(&commentInput)
	httputil.Response(ctx, status, nil)
}
func (m ChapterController) CreatePageComments(ctx *gin.Context) {
	commentInput, status := httputil.BindUriJson[dto2.PageCommentCreateInput](ctx)
	if status.IsError() {
		httputil.ErrorResponse(ctx, status)
		return
	}
	status = m.chapterService.CreatePageComment(&commentInput)
	httputil.Response(ctx, status, nil)
}
func (m ChapterController) FindChapterComments(ctx *gin.Context) {
	chapterId := ctx.Param("chapter_id")
	if len(chapterId) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}

	comments, status := m.chapterService.FindChapterComments(chapterId)
	httputil.Response(ctx, status, comments)
}
func (m ChapterController) FindPageComments(ctx *gin.Context) {
	pageId := ctx.Param("page_id")
	if len(pageId) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}
	comments, status := m.chapterService.FindPageComments(pageId)
	httputil.Response(ctx, status, comments)
}
func (m ChapterController) CreateChapter(ctx *gin.Context) {
	chapterInput, stat := httputil.BindUriJson[dto2.ChapterCreateInput](ctx)
	if stat.IsError() {
		httputil.ErrorResponse(ctx, stat)
		return
	}

	// Set translator / chapter creator
	claims, err := util.GetContextValue[*common.AccessTokenClaims](ctx, common.ClaimsKey)
	if err != nil {
		httputil.ErrorResponse(ctx, common.StatusError(status.AUTH_UNAUTHORIZED))
		return
	}
	chapterInput.TranslatorId = claims.UserId

	stat = m.chapterService.CreateChapter(&chapterInput)
	httputil.Response(ctx, stat, nil)
}
func (m ChapterController) DeleteChapter(ctx *gin.Context) {
	chapterId := ctx.Param("chapter_id")
	if len(chapterId) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}

	status := m.chapterService.DeleteChapter(chapterId)
	httputil.Response(ctx, status, nil)
}
func (m ChapterController) FindChapterPages(ctx *gin.Context) {
	chapterId := ctx.Param("chapter_id")
	if len(chapterId) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}

	pages, status := m.chapterService.FindChapterPages(chapterId)
	httputil.Response(ctx, status, pages)
}
func (m ChapterController) FindVolumeChapters(ctx *gin.Context) {
	//var volumeInput dto.VolumeDeleteInput
	//if err := ctx.BindUri(&volumeInput); err != nil {
	//	httputil.ErrorResponse(ctx, common.StatusError(common.BAD_PARAMETER_ERROR))
	//	return
	//}
	volumeId := ctx.Param("volume_id")
	if len(volumeId) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.BAD_PARAMETER_ERROR))
		return
	}

	chapters, status := m.chapterService.FindVolumeChapters(volumeId)
	httputil.Response(ctx, status, chapters)
}
