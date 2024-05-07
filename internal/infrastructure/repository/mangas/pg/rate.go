package pg

import (
  "context"
  "github.com/uptrace/bun"
  "manga-explorer/internal/domain/mangas"
  "manga-explorer/internal/domain/mangas/repository"
  "manga-explorer/internal/util"
  "time"
)

func NewMangaRate(db bun.IDB) repository.IRate {
  return &mangaRateRepository{db: db}
}

type mangaRateRepository struct {
  db bun.IDB
}

func (m mangaRateRepository) FindMangaRatings(mangaId string) ([]mangas.Rate, error) {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  var result []mangas.Rate
  err := m.db.NewSelect().
    Model(&result).
    Relation("User").
    Relation("Manga").
    Where("manga_id = ?", mangaId).
    Order("created_at").
    Scan(ctx)

  return util.CheckSliceResult(result, err).Unwrap()
}

func (m mangaRateRepository) FindRating(userId, mangaId string) (*mangas.Rate, error) {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  result := new(mangas.Rate)
  err := m.db.NewSelect().
    Model(result).
    Relation("User").
    Relation("Manga").
    Where("manga_id = ? AND user_id = ?", mangaId, userId).
    Scan(ctx)

  if err != nil {
    return nil, err
  }
  return result, nil
}

func (m mangaRateRepository) Upsert(rate *mangas.Rate) error {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  res, err := m.db.NewInsert().
    Model(rate).
    On("CONFLICT ON CONSTRAINT user_manga_idx DO UPDATE").
    Set("rate = EXCLUDED.rate").
    Set("updated_at = EXCLUDED.updated_at").
    Exec(ctx)

  return util.CheckSqlResult(res, err)
}
