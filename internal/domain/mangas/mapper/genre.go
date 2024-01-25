package mapper

import (
	"github.com/google/uuid"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
	"time"
)

func ToGenreResponse(genre *mangas.Genre) dto.GenreResponse {
	return dto.GenreResponse{
		Name: genre.Name,
	}
}

func MapGenreCreateInput(input dto.GenreCreateInput) mangas.Genre {
	return mangas.Genre{
		Id:        uuid.NewString(),
		Name:      input.Name,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

func MapGenreUpdateInput(input *dto.GenreUpdateInput) mangas.Genre {
	return mangas.Genre{
		Id:        uuid.NewString(),
		Name:      input.Name,
		UpdatedAt: time.Now(),
	}
}
