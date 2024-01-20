package repository

import "manga-explorer/internal/domain/mangas"

type IGenre interface {
	CreateGenre(genre *mangas.Genre) error
	DeleteGenreById(genreId string) error
	ListGenres() ([]mangas.Genre, error)
}
