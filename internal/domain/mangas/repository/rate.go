package repository

import (
	"manga-explorer/internal/domain/mangas"
)

type IRate interface {
	FindMangaRatings(mangaId string) ([]mangas.Rate, error)

	FindRating(userId, mangaId string) (*mangas.Rate, error)
	// Upsert each user rating some manga, it will be creating for the first time and updating for the rest
	Upsert(rate *mangas.Rate) error
}
