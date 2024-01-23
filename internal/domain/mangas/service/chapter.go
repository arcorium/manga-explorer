package service

import (
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/mangas/dto"
)

type IChapter interface {
	// CreateChapter Create new chapter which should be belonged to manga and volume
	CreateChapter(input *dto.ChapterCreateInput) status.Object
	// DeleteChapter Delete specific manga chapter by the id
	DeleteChapter(chapterId string) status.Object
	// EditChapter edit manga chapter
	EditChapter(input *dto.ChapterEditInput) status.Object
	// FindChapterPages Get manga chapter pages
	FindChapterPages(chapterId string) ([]dto.PageResponse, status.Object)
	// InsertChapterPage Uploads the image and set it as the page of manga chapter
	InsertChapterPage(input *dto.PageCreateInput) status.Object
	// CreateChapterComment Create new comment for manga chapter
	CreateChapterComment(input *dto.ChapterCommentCreateInput) status.Object
	// CreatePageComment Create new comment for chapter page
	CreatePageComment(input *dto.PageCommentCreateInput) status.Object
	// DeleteChapterPages Delete manga chapter pages based on the page numbers
	DeleteChapterPages(input *dto.PageDeleteInput) status.Object
	// FindVolumeChapters find all chapters in a volume
	FindVolumeChapters(volumeId string) ([]dto.ChapterResponse, status.Object)
	// FindChapterComments find all chapter comments
	FindChapterComments(chapterId string) ([]dto.CommentResponse, status.Object)
	// FindPageComments find all page comments
	FindPageComments(pageId string) ([]dto.CommentResponse, status.Object)
}
