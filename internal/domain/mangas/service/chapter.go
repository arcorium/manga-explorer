package service

import (
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/domain/mangas/dto"
)

type IChapter interface {
	// CreateChapter Create new chapter which should be belonged to manga and volume
	CreateChapter(input *dto.ChapterCreateInput) common.Status
	// DeleteChapter Delete specific manga chapter by the id
	DeleteChapter(chapterId string) common.Status
	// EditChapter edit manga chapter
	EditChapter(input *dto.ChapterEditInput) common.Status
	// FindChapterPages Get manga chapter pages
	FindChapterPages(chapterId string) ([]dto.PageResponse, common.Status)
	// InsertChapterPage Upload the image and set it as the page of manga chapter
	InsertChapterPage(input *dto.PageCreateInput) common.Status
	// CreateChapterComment Create new comment for manga chapter
	CreateChapterComment(input *dto.ChapterCommentCreateInput) common.Status
	// CreatePageComment Create new comment for chapter page
	CreatePageComment(input *dto.PageCommentCreateInput) common.Status
	// DeleteChapterPages Delete manga chapter pages based on the page numbers
	DeleteChapterPages(input *dto.PageDeleteInput) common.Status
	// FindVolumeChapters find all chapters in a volume
	FindVolumeChapters(volumeId string) ([]dto.ChapterResponse, common.Status)
	// FindChapterComments find all chapter comments
	FindChapterComments(chapterId string) ([]dto.CommentResponse, common.Status)
	// FindPageComments find all page comments
	FindPageComments(pageId string) ([]dto.CommentResponse, common.Status)
}
