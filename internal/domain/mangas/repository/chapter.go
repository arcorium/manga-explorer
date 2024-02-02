package repository

import (
	"manga-explorer/internal/domain/mangas"
	repo "manga-explorer/internal/infrastructure/repository"
)

type IChapter interface {
	CreateChapter(chapter *mangas.Chapter) error
	EditChapter(chapter *mangas.Chapter) error
	DeleteChapter(chapterId string) error
	FindChapter(id string) (*mangas.Chapter, error)
	FindVolumeDetails(volumeId string) (*mangas.Volume, error)
	FindPagesDetails(chapterId string, pages []uint16) ([]mangas.Page, error)
	DeleteChapterPages(chapterId string, pages []uint16) error
	InsertChapterPages(pages []mangas.Page) error
	InsertChapterHistories(history *mangas.ChapterHistory) error
	FindMangaChapterHistories(userId string, mangaId string, pagedQuery repo.QueryParameter) (repo.PagedQueryResult[[]mangas.Chapter], error)
}
