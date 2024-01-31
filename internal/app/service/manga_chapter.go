package service

import (
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/mangas"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/infrastructure/file"
	fileService "manga-explorer/internal/infrastructure/file/service"
	"manga-explorer/internal/util/containers"
	"manga-explorer/internal/util/opt"
	"mime/multipart"
)

func NewChapterService(chapterRepo repository.IChapter, commentRepo repository.IComment) service.IChapter {
	return &mangaChapterService{
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

func (m mangaChapterService) FindChapterPages(chapterId string) ([]mangaDto.PageResponse, status.Object) {
	pages, err := m.chapterRepo.FindChapterPages(chapterId)
	pageResponses := containers.CastSlicePtr1(pages, m.fileService, mapper.ToPageResponse)
	return pageResponses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaChapterService) CreateChapter(input *mangaDto.ChapterCreateInput) status.Object {
	chapter := mapper.MapChapterCreateInput(input)
	err := m.chapterRepo.CreateChapter(&chapter)
	return status.ConditionalRepository(err, status.CREATED, opt.New(status.CHAPTER_ALREADY_EXIST))
}

func (m mangaChapterService) InsertChapterPage(input *mangaDto.PageCreateInput) status.Object {
	fileHeaders := containers.CastSlicePtr(input.Pages, func(current *mangaDto.InternalPage) multipart.FileHeader {
		return *current.Image
	})
	filenames, stat := m.fileService.Uploads(file.MangaAsset, fileHeaders)
	if stat.IsError() {
		return stat
	}

	pages := mapper.MapPageCreateInput(input, filenames)
	if pages == nil {
		return status.InternalError()
	}

	err := m.chapterRepo.InsertChapterPages(pages)
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
	err := m.chapterRepo.DeleteChapterPages(input.ChapterId, input.Pages)
	return status.ConditionalRepository(err, status.DELETED, opt.New(status.PAGE_NOT_FOUND))
}

func (m mangaChapterService) FindVolumeChapters(volumeId string) ([]mangaDto.ChapterResponse, status.Object) {
	chapters, err := m.chapterRepo.FindVolumeChapters(volumeId)
	chapterResponses := containers.CastSlicePtr1(chapters, m.fileService, mapper.ToChapterResponse)
	return chapterResponses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS)) // Make empty response as success
}

func (m mangaChapterService) FindChapterComments(chapterId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindChapterComments(chapterId)
	responses := mapper.ToCommentResponse2(comments)
	return responses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaChapterService) FindPageComments(pageId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindPageComments(pageId)
	responses := mapper.ToCommentResponse2(comments)
	return responses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}
