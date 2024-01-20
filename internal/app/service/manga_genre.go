package service

import (
	"manga-explorer/internal/app/common"
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

func (m mangaGenreService) CreateGenre(input dto.GenreCreateInput) common.Status {
	genre := mapper.MapGenreCreateInput(input)
	err := m.genreRepo.CreateGenre(&genre)
	return common.NewRepositoryStatus(err, status.SUCCESS_CREATED)
}

func (m mangaGenreService) DeleteGenre(genreId string) common.Status {
	err := m.genreRepo.DeleteGenreById(genreId)
	return common.NewRepositoryStatus(err)
}

func (m mangaGenreService) ListGenre() ([]dto.GenreResponse, common.Status) {
	genres, err := m.genreRepo.ListGenres()
	if err != nil {
		return nil, common.NewRepositoryStatus(err)
	}
	genreResponses := containers.CastSlicePtr(genres, mapper.ToGenreResponse)
	return genreResponses, common.StatusSuccess()
}
