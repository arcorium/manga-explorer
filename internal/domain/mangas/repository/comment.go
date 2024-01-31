package repository

import "manga-explorer/internal/domain/mangas"

type IComment interface {
	FindMangaComments(mangaId string) ([]mangas.Comment, error)
	FindChapterComments(chapterId string) ([]mangas.Comment, error)
	FindPageComments(pageId string) ([]mangas.Comment, error)
	FindComment(id string) (*mangas.Comment, error)
	CreateComment(comment *mangas.Comment) error
	EditComment(comment *mangas.Comment) error
	DeleteComment(commentId string) error
}
