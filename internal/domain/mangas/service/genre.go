package service

import (
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/domain/mangas/dto"
)

type IGenre interface {
	// CreateGenre create new genre
	CreateGenre(input dto.GenreCreateInput) common.Status
	// DeleteGenre delete genre by the id
	DeleteGenre(genreId string) common.Status
	// ListGenre get all available genres
	ListGenre() ([]dto.GenreResponse, common.Status)
}
