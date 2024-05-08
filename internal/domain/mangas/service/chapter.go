package service

import (
	dto2 "manga-explorer/internal/common/dto"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/util/opt"
)

type IChapter interface {
	// CreateChapter Upsert new chapter which should be belonged to manga and volume
	CreateChapter(input *dto.ChapterCreateInput) status.Object
	// DeleteChapter Delete specific manga chapter by the id
	DeleteChapter(chapterId string) status.Object
	// EditChapter edit manga chapter
	EditChapter(input *dto.ChapterEditInput) status.Object
	// FindChapterDetails Get manga chapter pages
	FindMangaChapterHistories(input *dto.MangaChapterHistoriesFindInput) ([]dto.ChapterResponse, *dto2.ResponsePage, status.Object)
	FindChapterDetails(chapterId string, userId opt.Optional[string]) (dto.ChapterResponse, status.Object)
	// InsertChapterPage Uploads the image and set it as the page of manga chapter, it will return pages that failed to be inserted
	InsertChapterPage(input *dto.PageCreateInput) (status.Object, []uint16)
	// CreateChapterComment Upsert new comment for manga chapter
	CreateChapterComment(input *dto.ChapterCommentCreateInput) status.Object
	// CreatePageComment Upsert new comment for chapter page
	CreatePageComment(input *dto.PageCommentCreateInput) status.Object
	// DeleteChapterPages Delete manga chapter pages based on the page numbers
	DeleteChapterPages(input *dto.PageDeleteInput) status.Object
	// FindVolumeDetails find all chapters in a volume
	FindVolumeDetails(volumeId string) (dto.VolumeResponse, status.Object)
	// FindChapterComments find all chapter comments
	FindChapterComments(chapterId string) ([]dto.CommentResponse, status.Object)
	// FindPageComments find all page comments
	FindPageComments(pageId string) ([]dto.CommentResponse, status.Object)
}
