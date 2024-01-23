package service

import (
	"manga-explorer/internal/app/common/status"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/infrastructure/file"
	fileService "manga-explorer/internal/infrastructure/file/service"
	"manga-explorer/internal/util/containers"
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
	return status.ConditionalRepository(err, status.SUCCESS)
}

func (m mangaChapterService) FindChapterPages(chapterId string) ([]mangaDto.PageResponse, status.Object) {
	pages, err := m.chapterRepo.FindChapterPages(chapterId)
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	pageResponses := containers.CastSlicePtr1(pages, m.fileService, mapper.ToPageResponse)
	return pageResponses, status.Success()
}

func (m mangaChapterService) CreateChapter(input *mangaDto.ChapterCreateInput) status.Object {
	chapter := mapper.MapChapterCreateInput(input)
	err := m.chapterRepo.CreateChapter(&chapter)
	return status.ConditionalRepository(err, status.CREATED)
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
	return status.ConditionalRepository(err, status.UPDATED)
}

func (m mangaChapterService) EditChapter(input *mangaDto.ChapterEditInput) status.Object {
	chapter := mapper.MapChapterEditInput(input)
	err := m.chapterRepo.EditChapter(&chapter)
	return status.ConditionalRepository(err, status.UPDATED)
}

func (m mangaChapterService) CreateChapterComment(input *mangaDto.ChapterCommentCreateInput) status.Object {
	comment := mapper.MapChapterCommentCreateInput(input)
	if input.HasParent() {
		parent, err := m.commentRepo.FindComment(input.ParentId)
		if err != nil {
			return status.RepositoryError(err)
		}

		// Validate
		if !comment.ValidateAsReply(parent) {
			return status.Error(status.BAD_BODY_REQUEST_ERROR, "Couldn't reply to such comment")
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return status.ConditionalRepository(err, status.CREATED)
}

func (m mangaChapterService) CreatePageComment(input *mangaDto.PageCommentCreateInput) status.Object {
	comment := mapper.MapPageCommentCreateInput(input)
	if input.HasParent() {
		parent, err := m.commentRepo.FindComment(input.ParentId)
		if err != nil {
			return status.RepositoryError(err)
		}

		// Validate
		if !comment.ValidateAsReply(parent) {
			return status.Error(status.BAD_BODY_REQUEST_ERROR, "Couldn't reply to such comment")
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return status.ConditionalRepository(err, status.CREATED)
}

func (m mangaChapterService) DeleteChapterPages(input *mangaDto.PageDeleteInput) status.Object {
	err := m.chapterRepo.DeleteChapterPages(input.ChapterId, input.Pages)
	return status.ConditionalRepository(err, status.DELETED)
}

func (m mangaChapterService) FindVolumeChapters(volumeId string) ([]mangaDto.ChapterResponse, status.Object) {
	chapters, err := m.chapterRepo.FindVolumeChapters(volumeId)
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	chapterResponses := containers.CastSlicePtr1(chapters, m.fileService, mapper.ToChapterResponse)
	return chapterResponses, status.Success()
}

func (m mangaChapterService) FindChapterComments(chapterId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindChapterComments(chapterId)
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	commentResponses := containers.CastSlicePtr(comments, mapper.ToCommentResponse)
	return commentResponses, status.Success()
}

func (m mangaChapterService) FindPageComments(pageId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindPageComments(pageId)
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	commentResponses := containers.CastSlicePtr(comments, mapper.ToCommentResponse)
	return commentResponses, status.Success()
}
