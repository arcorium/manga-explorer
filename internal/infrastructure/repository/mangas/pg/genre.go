package pg

import (
	"context"
	"github.com/uptrace/bun"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/util"
	"time"
)

func NewMangaGenre(db bun.IDB) repository.IGenre {
	return &mangaGenreRepository{db: db}
}

type mangaGenreRepository struct {
	db bun.IDB
}

func (m mangaGenreRepository) CreateGenre(genre *mangas.Genre) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewInsert().
		Model(genre).
		Returning("NULL").
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (m mangaGenreRepository) UpdateGenre(genre *mangas.Genre) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewUpdate().
		Model(genre).
		WherePK().
		OmitZero().
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (m mangaGenreRepository) DeleteGenreById(genreId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewDelete().
		Model((*mangas.Genre)(nil)).
		Where("id = ?", genreId).
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (m mangaGenreRepository) ListGenres() ([]mangas.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.Genre
	err := m.db.NewSelect().
		Model(&result).
		Scan(ctx)

	return util.CheckSliceResult(result, err).Unwrap()
}
