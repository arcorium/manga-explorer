package repository

import "manga-explorer/internal/domain/mangas"

type IChapter interface {
	CreateChapter(chapter *mangas.Chapter) error
	EditChapter(chapter *mangas.Chapter) error
	DeleteChapter(chapterId string) error
	FindChapter(id string) (*mangas.Chapter, error)
	FindChapterPages(chapterId string) ([]mangas.Page, error)
	FindVolumeChapters(volumeId string) ([]mangas.Chapter, error)
	DeleteChapterPages(chapterId string, pages []uint16) error
	InsertChapterPages(pages []mangas.Page) error
}
