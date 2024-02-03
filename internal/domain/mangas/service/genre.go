package service

import (
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/mangas/dto"
)

type IGenre interface {
	// CreateGenre create new genre
	CreateGenre(input dto.GenreCreateInput) status.Object
	// DeleteGenre delete genre by the id
	DeleteGenre(genreId string) status.Object
	// UpdateGenre update genre name
	UpdateGenre(input *dto.GenreEditInput) status.Object
	// ListGenre get all available genres
	ListGenre() ([]dto.GenreResponse, status.Object)
}
