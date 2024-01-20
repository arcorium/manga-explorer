package mangas

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/util/httputil"
	"manga-explorer/internal/util/httputil/resp"
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
	var input dto.PageCreateInput
	stat, fieldsErr := httputil.BindUriMultipartForm(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}
	stat = m.chapterService.InsertChapterPage(&input)
	resp.Conditional(ctx, stat, nil, nil)
}
func (m ChapterController) DeleteChapterPage(ctx *gin.Context) {
	var pageInput dto.PageDeleteInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &pageInput)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
	}

	stat = m.chapterService.DeleteChapterPages(&pageInput)
	resp.Conditional(ctx, stat, nil, nil)
}
func (m ChapterController) EditChapter(ctx *gin.Context) {
	var editInput dto.ChapterEditInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &editInput)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.chapterService.EditChapter(&editInput)
	resp.Conditional(ctx, stat, nil, nil)
}
func (m ChapterController) CreateChapterComments(ctx *gin.Context) {
	var input dto.ChapterCommentCreateInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}

	stat = m.chapterService.CreateChapterComment(&input)
	resp.Conditional(ctx, stat, nil, nil)
}
func (m ChapterController) CreatePageComments(ctx *gin.Context) {
	var input dto.PageCommentCreateInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &input)
	if stat.IsError() {
		resp.ErrorDetailed(ctx, stat, fieldsErr)
		return
	}
	stat = m.chapterService.CreatePageComment(&input)
	resp.Conditional(ctx, stat, nil, nil)
}
func (m ChapterController) FindChapterComments(ctx *gin.Context) {
	chapterId := ctx.Param("chapter_id")
	if len(chapterId) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("chapter_id"))
		return
	}

	comments, stat := m.chapterService.FindChapterComments(chapterId)
	resp.Conditional(ctx, stat, comments, nil)
}
func (m ChapterController) FindPageComments(ctx *gin.Context) {
	pageId := ctx.Param("page_id")
	if len(pageId) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("page_id"))
		return
	}
	comments, stat := m.chapterService.FindPageComments(pageId)
	resp.Conditional(ctx, stat, comments, nil)
}
func (m ChapterController) CreateChapter(ctx *gin.Context) {
	// TODO: Add support for creating chapter and inserting pages so each chapter should not have empty pages
	var input dto.ChapterCreateInput
	stat, fieldsErr := httputil.BindUriJson(ctx, &input)
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
func (m ChapterController) DeleteChapter(ctx *gin.Context) {
	chapterId := ctx.Param("chapter_id")
	if len(chapterId) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("chapter_id"))
		return
	}

	stat := m.chapterService.DeleteChapter(chapterId)
	resp.Conditional(ctx, stat, nil, nil)
}
func (m ChapterController) FindChapterPages(ctx *gin.Context) {
	chapterId := ctx.Param("chapter_id")
	if len(chapterId) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("chapter_id"))
		return
	}

	pages, stat := m.chapterService.FindChapterPages(chapterId)
	resp.Conditional(ctx, stat, pages, nil)
}
func (m ChapterController) FindVolumeChapters(ctx *gin.Context) {
	volumeId := ctx.Param("volume_id")
	if len(volumeId) == 0 {
		resp.ErrorDetailed(ctx, status.Error(status.BAD_PARAMETER_ERROR), common.NewNotPresentParameter("volume_id"))
		return
	}

	chapters, stat := m.chapterService.FindVolumeChapters(volumeId)
	resp.Conditional(ctx, stat, chapters, nil)
}
