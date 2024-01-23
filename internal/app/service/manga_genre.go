package service

import (
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/util/containers"
)

func NewGenreService(genreRepo repository.IGenre) service.IGenre {
	return &mangaGenreService{genreRepo: genreRepo}
}

type mangaGenreService struct {
	genreRepo repository.IGenre
}

func (m mangaGenreService) CreateGenre(input dto.GenreCreateInput) status.Object {
	genre := mapper.MapGenreCreateInput(input)
	err := m.genreRepo.CreateGenre(&genre)
	return status.ConditionalRepository(err, status.CREATED)
}

func (m mangaGenreService) DeleteGenre(genreId string) status.Object {
	err := m.genreRepo.DeleteGenreById(genreId)
	return status.ConditionalRepository(err, status.DELETED)
}

func (m mangaGenreService) ListGenre() ([]dto.GenreResponse, status.Object) {
	genres, err := m.genreRepo.ListGenres()
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	genreResponses := containers.CastSlicePtr(genres, mapper.ToGenreResponse)
	return genreResponses, status.Success()
}
