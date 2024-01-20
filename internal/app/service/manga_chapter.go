package service

import (
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/mangas"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/infrastructure/file"
	fileService "manga-explorer/internal/infrastructure/file/service"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/containers"
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
	return status.FromRepository(err)
}

func (m mangaChapterService) FindChapterPages(chapterId string) ([]mangaDto.PageResponse, status.Object) {
	pages, err := m.chapterRepo.FindChapterPages(chapterId)
	stat := status.FromRepository(err)
	if stat.IsError() {
		return nil, stat
	}
	pageResponses := containers.CastSlicePtr(pages, mapper.ToPageResponse)
	return pageResponses, stat
}

func (m mangaChapterService) CreateChapter(input *mangaDto.ChapterCreateInput) status.Object {
	chapter := mapper.MapChapterCreateInput(input)
	err := m.chapterRepo.CreateChapter(&chapter)
	return status.FromRepository(err, status.CREATED)
}

func (m mangaChapterService) InsertChapterPage(input *mangaDto.PageCreateInput) status.Object {
	page := mapper.MapPageCreateInput(input)
	format, err := util.GetFileFormat(input.Image.Filename)
	if err != nil {
		return status.Error(status.INTERNAL_SERVER_ERROR)
	}
	page.ImageURL = m.fileService.GetURL(file.ImageAsset, page.Id, format)

	// Read bytes
	fl, err := input.Image.Open()
	if err != nil {
		return status.Error(status.INTERNAL_SERVER_ERROR)
	}

	var bytes []byte
	_, err = fl.Read(bytes)
	if err != nil {
		return status.Error(status.INTERNAL_SERVER_ERROR)
	}
	// Upload image
	m.fileService.Upload(page.ImageURL, bytes)

	err = m.chapterRepo.InsertChapterPages([]mangas.Page{page})
	return status.FromRepository(err, status.UPDATED)
}

func (m mangaChapterService) EditChapter(input *mangaDto.ChapterEditInput) status.Object {
	chapter := mapper.MapChapterEditInput(input)
	err := m.chapterRepo.EditChapter(&chapter)
	return status.FromRepository(err, status.UPDATED)
}

func (m mangaChapterService) CreateChapterComment(input *mangaDto.ChapterCommentCreateInput) status.Object {
	comment := mapper.MapChapterCommentCreateInput(input)
	if input.HasParent() {
		parent, err := m.commentRepo.FindComment(input.ParentId)
		stat := status.FromRepository(err)
		if stat.IsError() {
			return status.Error(status.INTERNAL_SERVER_ERROR)
		}

		// Validate
		if !comment.ValidateAsReply(parent) {
			return status.Error(status.BAD_BODY_REQUEST_ERROR, "Couldn't reply to such comment")
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return status.FromRepository(err, status.CREATED)
}

func (m mangaChapterService) CreatePageComment(input *mangaDto.PageCommentCreateInput) status.Object {
	comment := mapper.MapPageCommentCreateInput(input)
	if input.HasParent() {
		parent, err := m.commentRepo.FindComment(input.ParentId)
		stat := status.FromRepository(err)
		if stat.IsError() {
			return status.Error(status.INTERNAL_SERVER_ERROR)
		}

		// Validate
		if !comment.ValidateAsReply(parent) {
			return status.Error(status.BAD_BODY_REQUEST_ERROR, "Couldn't reply to such comment")
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return status.FromRepository(err, status.CREATED)
}

func (m mangaChapterService) DeleteChapterPages(input *mangaDto.PageDeleteInput) status.Object {
	err := m.chapterRepo.DeleteChapterPages(input.ChapterId, input.Pages)
	return status.FromRepository(err, status.DELETED)
}

func (m mangaChapterService) FindVolumeChapters(volumeId string) ([]mangaDto.ChapterResponse, status.Object) {
	chapters, err := m.chapterRepo.FindVolumeChapters(volumeId)
	stat := status.FromRepository(err)
	if stat.IsError() {
		return nil, stat
	}
	chapterResponses := containers.CastSlicePtr(chapters, mapper.ToChapterResponse)
	return chapterResponses, stat
}

func (m mangaChapterService) FindChapterComments(chapterId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindChapterComments(chapterId)
	stat := status.FromRepository(err)
	if stat.IsError() {
		return nil, stat
	}
	commentResponses := containers.CastSlicePtr(comments, mapper.ToCommentResponse)
	return commentResponses, stat
}

func (m mangaChapterService) FindPageComments(pageId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindPageComments(pageId)
	stat := status.FromRepository(err)
	if stat.IsError() {
		return nil, stat
	}
	commentResponses := containers.CastSlicePtr(comments, mapper.ToCommentResponse)
	return commentResponses, stat
}
