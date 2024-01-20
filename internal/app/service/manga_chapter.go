package service

import (
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/app/service/utility/file"
	"manga-explorer/internal/domain/mangas"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
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
	fileService file.IService

	chapterRepo repository.IChapter
	commentRepo repository.IComment
}

func (m mangaChapterService) DeleteChapter(chapterId string) common.Status {
	err := m.chapterRepo.DeleteChapter(chapterId)
	return common.NewRepositoryStatus(err)
}

func (m mangaChapterService) FindChapterPages(chapterId string) ([]mangaDto.PageResponse, common.Status) {
	pages, err := m.chapterRepo.FindChapterPages(chapterId)
	if err != nil {
		return nil, common.NewRepositoryStatus(err)
	}
	pageResponses := containers.CastSlicePtr(pages, mapper.ToPageResponse)
	return pageResponses, common.StatusSuccess()
}

func (m mangaChapterService) CreateChapter(input *mangaDto.ChapterCreateInput) common.Status {
	chapter := mapper.MapChapterCreateInput(input)
	err := m.chapterRepo.CreateChapter(&chapter)
	return common.NewRepositoryStatus(err, status.SUCCESS_CREATED)
}

func (m mangaChapterService) InsertChapterPage(input *mangaDto.PageCreateInput) common.Status {
	page := mapper.MapPageCreateInput(input)
	format, err := util.GetFileFormat(input.Image.Filename)
	if err != nil {
		return common.StatusError(status.INTERNAL_SERVER_ERROR)
	}
	page.ImageURL = m.fileService.GetURL(file.ImageAsset, page.Id, format)

	// Read bytes
	fl, err := input.Image.Open()
	if err != nil {
		return common.StatusError(status.INTERNAL_SERVER_ERROR)
	}

	var bytes []byte
	_, err = fl.Read(bytes)
	if err != nil {
		return common.StatusError(status.INTERNAL_SERVER_ERROR)
	}
	// Upload image
	m.fileService.Upload(page.ImageURL, bytes)

	err = m.chapterRepo.InsertChapterPages([]mangas.Page{page})
	return common.NewRepositoryStatus(err, status.SUCCESS_CREATED)
}

func (m mangaChapterService) EditChapter(input *mangaDto.ChapterEditInput) common.Status {
	chapter := mapper.MapChapterEditInput(input)
	err := m.chapterRepo.EditChapter(&chapter)
	return common.NewRepositoryStatus(err)
}

func (m mangaChapterService) CreateChapterComment(input *mangaDto.ChapterCommentCreateInput) common.Status {
	comment := mapper.MapChapterCommentCreateInput(input)
	if input.HasParent() {
		parent, err := m.commentRepo.FindComment(input.ParentId)
		if err != nil {
			return common.StatusError(status.INTERNAL_SERVER_ERROR)
		}

		// Validate
		if !comment.ValidateAsReply(parent) {
			return common.StatusError(status.BAD_BODY_REQUEST_ERROR)
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return common.NewRepositoryStatus(err, status.SUCCESS_CREATED)
}

func (m mangaChapterService) CreatePageComment(input *mangaDto.PageCommentCreateInput) common.Status {
	comment := mapper.MapPageCommentCreateInput(input)
	if input.HasParent() {
		parent, err := m.commentRepo.FindComment(input.ParentId)
		if err != nil {
			return common.StatusError(status.INTERNAL_SERVER_ERROR)
		}

		// Validate
		if !comment.ValidateAsReply(parent) {
			return common.StatusError(status.BAD_BODY_REQUEST_ERROR)
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return common.NewRepositoryStatus(err, status.SUCCESS_CREATED)
}

func (m mangaChapterService) DeleteChapterPages(input *mangaDto.PageDeleteInput) common.Status {
	err := m.chapterRepo.DeleteChapterPages(input.ChapterId, input.Pages)
	return common.NewRepositoryStatus(err)
}

func (m mangaChapterService) FindVolumeChapters(volumeId string) ([]mangaDto.ChapterResponse, common.Status) {
	chapters, err := m.chapterRepo.FindVolumeChapters(volumeId)
	if err != nil {
		return nil, common.NewRepositoryStatus(err)
	}
	chapterResponses := containers.CastSlicePtr(chapters, mapper.ToChapterResponse)
	return chapterResponses, common.StatusSuccess()
}

func (m mangaChapterService) FindChapterComments(chapterId string) ([]mangaDto.CommentResponse, common.Status) {
	comments, err := m.commentRepo.FindChapterComments(chapterId)
	if err != nil {
		return nil, common.NewRepositoryStatus(err)
	}
	commentResponses := containers.CastSlicePtr(comments, mapper.ToCommentResponse)
	return commentResponses, common.StatusSuccess()
}

func (m mangaChapterService) FindPageComments(pageId string) ([]mangaDto.CommentResponse, common.Status) {
	comments, err := m.commentRepo.FindPageComments(pageId)
	if err != nil {
		return nil, common.NewRepositoryStatus(err)
	}
	commentResponses := containers.CastSlicePtr(comments, mapper.ToCommentResponse)
	return commentResponses, common.StatusSuccess()
}
