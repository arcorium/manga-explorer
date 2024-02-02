package repository

import (
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/infrastructure/repository"
)

type IManga interface {
	CreateManga(manga *mangas.Manga, genres []mangas.MangaGenre) error
	EditManga(manga *mangas.Manga) error
	PatchManga(manga *mangas.Manga) error
	EditMangaGenres(additional, removes []mangas.MangaGenre) error
	FindMinimalMangaById(id string) (*mangas.Manga, error)
	FindMangasById(ids ...string) ([]mangas.Manga, error)
	// FindMangasByFilter Get manga based on the filter specified, set limit and offset both to 0 to get all the mangas
	FindMangasByFilter(filter *mangas.SearchFilter, pagedQuery repository.QueryParameter) (repository.PagedQueryResult[[]mangas.Manga], error)
	// FindRandomMangas Get manga which will be returning different manga for each call, set limit to 0 to get all the mangas
	FindRandomMangas(limit uint64) ([]mangas.Manga, error)
	FindMangaHistories(userId string, pagedQuery repository.QueryParameter) (repository.PagedQueryResult[[]mangas.MangaHistory], error)
	// FindMangaFavorites Find favorites mangas by userId, returning favorites mangas and total favorites mangas on user
	FindMangaFavorites(userId string, pagedQuery repository.QueryParameter) (repository.PagedQueryResult[[]mangas.MangaFavorite], error)
	InsertMangaFavorite(favorite *mangas.MangaFavorite) error
	RemoveMangaFavorite(favorite *mangas.MangaFavorite) error
	// ListMangas Get all manga based on the offset and limit, set limit and offset both to 0 to get all the mangas
	ListMangas(parameter repository.QueryParameter) (repository.PagedQueryResult[[]mangas.Manga], error)
	CreateVolume(volume *mangas.Volume) error
	DeleteVolume(mangaId string, volume uint32) error
}
