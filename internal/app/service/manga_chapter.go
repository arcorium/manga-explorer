package service

import (
	dto2 "manga-explorer/internal/common/dto"
	mapper2 "manga-explorer/internal/common/mapper"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/mangas"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/infrastructure/file"
	fileService "manga-explorer/internal/infrastructure/file/service"
	repo "manga-explorer/internal/infrastructure/repository"
	"manga-explorer/internal/util/containers"
	"manga-explorer/internal/util/opt"
)

func NewChapterService(fileService fileService.IFile, chapterRepo repository.IChapter, commentRepo repository.IComment) service.IChapter {
	return &mangaChapterService{
		fileService: fileService,
		chapterRepo: chapterRepo,
		commentRepo: commentRepo,
	}
}

type mangaChapterService struct {
	fileService fileService.IFile

	chapterRepo repository.IChapter
	commentRepo repository.IComment
}

func (m mangaChapterService) DeleteChapter(chapterId string) status.Object {
	err := m.chapterRepo.DeleteChapter(chapterId)
	return status.ConditionalRepository(err, status.SUCCESS, opt.New(status.CHAPTER_NOT_FOUND))
}

func (m mangaChapterService) FindChapterDetails(chapterId string, userId opt.Optional[string]) (mangaDto.ChapterResponse, status.Object) {
	chapter, err := m.chapterRepo.FindChapter(chapterId)
	responses := mapper.ToChapterResponse(chapter, m.fileService)
	// Add chapter history
	if userId.HasValue() && err == nil {
		chapterHistory := mangas.NewChapterHistory(*userId.Value(), chapterId, opt.NullTime)
		err = m.chapterRepo.InsertChapterHistories(&chapterHistory)
		if err != nil {
			return mangaDto.ChapterResponse{}, status.RepositoryError(err, opt.New(status.CHAPTER_UPDATE_FAILED))
		}
	}
	return responses, status.ConditionalRepositoryE(err, status.SUCCESS, opt.New(status.SUCCESS), opt.New(status.OBJECT_NOT_FOUND))
}

func (m mangaChapterService) FindMangaChapterHistories(input *mangaDto.MangaChapterHistoriesFindInput) ([]mangaDto.ChapterResponse, *dto2.ResponsePage, status.Object) {
	chapterHistories, err := m.chapterRepo.FindMangaChapterHistories(input.UserId, input.MangaId, repo.QueryParameter{input.Offset(), input.Element})
	page := mapper2.NewResponsePage(chapterHistories.Data, chapterHistories.Total, &input.PagedQueryInput)
	responses := containers.CastSlicePtr(chapterHistories.Data, mapper.ToMinimalChapterResponse)
	return responses, &page, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaChapterService) CreateChapter(input *mangaDto.ChapterCreateInput) status.Object {
	chapter := mapper.MapChapterCreateInput(input)
	err := m.chapterRepo.CreateChapter(&chapter)
	return status.ConditionalRepository(err, status.CREATED, opt.New(status.CHAPTER_ALREADY_EXIST))
}

//func (m mangaChapterService) InsertChapterPage(input *mangaDto.PageCreateInput) status.Object {
//	fileHeaders := containers.CastSlicePtr(input.Page, func(current *mangaDto.InternalPage) multipart.FileHeader {
//		return *current.Image
//	})
//	filenames, stat := m.fileService.Uploads(file.MangaAsset, fileHeaders)
//	if stat.IsError() {
//		return stat
//	}
//
//	pages := mapper.MapPageCreateInput(input, filenames)
//	if pages == nil {
//		return status.InternalError()
//	}
//
//	err := m.chapterRepo.InsertChapterPages(pages)
//	return status.ConditionalRepository(err, status.UPDATED, opt.New(status.PAGE_INSERT_FAILED))
//}

func (m mangaChapterService) InsertChapterPage(input *mangaDto.PageCreateInput) status.Object {
	filename, stat := m.fileService.Upload(file.MangaAsset, input.Page.Image)
	if stat.IsError() {
		return stat
	}

	pages := mapper.MapPageCreateInput(input, filename)

	err := m.chapterRepo.InsertChapterPages([]mangas.Page{pages})
	if err != nil {
		m.fileService.Delete(file.MangaAsset, filename)
	}
	return status.ConditionalRepository(err, status.UPDATED, opt.New(status.PAGE_INSERT_FAILED))
}

func (m mangaChapterService) EditChapter(input *mangaDto.ChapterEditInput) status.Object {
	chapter := mapper.MapChapterEditInput(input)
	err := m.chapterRepo.EditChapter(&chapter)
	return status.ConditionalRepository(err, status.UPDATED, opt.New(status.CHAPTER_UPDATE_FAILED))
}

func (m mangaChapterService) CreateChapterComment(input *mangaDto.ChapterCommentCreateInput) status.Object {
	comment := mapper.MapChapterCommentCreateInput(input)
	if input.HasParent() {
		if stat := m.validateReplyComment(input.ParentId, &comment); stat.IsError() {
			return stat
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return status.ConditionalRepository(err, status.CREATED, opt.New(status.COMMENT_CREATE_FAILED))
}

func (m mangaChapterService) CreatePageComment(input *mangaDto.PageCommentCreateInput) status.Object {
	comment := mapper.MapPageCommentCreateInput(input)
	if input.HasParent() {
		if stat := m.validateReplyComment(input.ParentId, &comment); stat.IsError() {
			return stat
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return status.ConditionalRepository(err, status.CREATED, opt.New(status.COMMENT_CREATE_FAILED))
}

func (m mangaChapterService) validateReplyComment(parentId string, comment *mangas.Comment) status.Object {
	parent, err := m.commentRepo.FindComment(parentId)
	if err != nil {
		return status.RepositoryError(err, opt.New(status.COMMENT_PARENT_NOT_FOUND))
	}

	// Response
	if !comment.ValidateAsReply(parent) {
		return status.Error(status.COMMENT_PARENT_DIFFERENT_SCOPE)
	}
	return status.Success()
}

func (m mangaChapterService) DeleteChapterPages(input *mangaDto.PageDeleteInput) status.Object {
	// Get pages details for deleting the page images
	pagesDetails, err := m.chapterRepo.FindPagesDetails(input.ChapterId, input.Pages)
	if err != nil {
		return status.RepositoryError(err, opt.New(status.PAGE_NOT_FOUND))
	}

	// Delete metadata
	err = m.chapterRepo.DeleteChapterPages(input.ChapterId, input.Pages)
	if err != nil {
		return status.RepositoryError(err, opt.New(status.PAGE_NOT_FOUND))
	}

	// Delete files
	for _, val := range pagesDetails {
		// Skip bad value
		if len(val.ImageURL) > 0 {
			m.fileService.Delete(file.MangaAsset, val.ImageURL)
		}
	}

	return status.ConditionalRepository(err, status.DELETED, opt.New(status.PAGE_NOT_FOUND))
}

func (m mangaChapterService) FindVolumeDetails(volumeId string) (mangaDto.VolumeResponse, status.Object) {
	chapters, err := m.chapterRepo.FindVolumeDetails(volumeId)
	if err != nil {
		return mangaDto.VolumeResponse{}, status.RepositoryError(err, opt.New(status.OBJECT_NOT_FOUND))
	}
	response := mapper.ToVolumeResponse(chapters, m.fileService)
	return response, status.Success()
}

func (m mangaChapterService) FindChapterComments(chapterId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindChapterComments(chapterId)
	responses := mapper.ToCommentsResponse(comments)
	return responses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaChapterService) FindPageComments(pageId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindPageComments(pageId)
	responses := mapper.ToCommentsResponse(comments)
	return responses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}
