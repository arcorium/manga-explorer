package service

import (
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/util/containers"
	"manga-explorer/internal/util/opt"
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
	return status.ConditionalRepository(err, status.CREATED, opt.New(status.GENRE_ALREADY_EXIST))
}

func (m mangaGenreService) DeleteGenre(genreId string) status.Object {
	err := m.genreRepo.DeleteGenreById(genreId)
	return status.ConditionalRepository(err, status.DELETED, opt.New(status.GENRE_NOT_FOUND))
}

func (m mangaGenreService) UpdateGenre(input *dto.GenreEditInput) status.Object {
	updatedGenre := mapper.MapGenreUpdateInput(input)
	err := m.genreRepo.UpdateGenre(&updatedGenre)
	return status.ConditionalRepository(err, status.UPDATED, opt.New(status.GENRE_NOT_FOUND))
}

func (m mangaGenreService) ListGenre() ([]dto.GenreResponse, status.Object) {
	genres, err := m.genreRepo.ListGenres()
	genreResponses := containers.CastSlicePtr(genres, mapper.ToGenreResponse)
	return genreResponses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}
