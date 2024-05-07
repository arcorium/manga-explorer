package pg

import (
  "context"
  "github.com/uptrace/bun"
  "manga-explorer/internal/common"
  "manga-explorer/internal/domain/mangas"
  "manga-explorer/internal/domain/mangas/repository"
  "manga-explorer/internal/util"
  "time"
)

func NewTranslationRepository(db bun.IDB) repository.ITranslation {
  return &translationRepository{db: db}
}

type translationRepository struct {
  db bun.IDB
}

func (t translationRepository) Create(translation []mangas.Translation) error {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  res, err := t.db.NewInsert().
    Model(&translation).
    Returning("NULL").
    Exec(ctx)
  return util.CheckSqlResult(res, err)
}

func (t translationRepository) FindByMangaId(mangaId string) ([]mangas.Translation, error) {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  var result []mangas.Translation
  err := t.db.NewSelect().
    Model(&result).
    Relation("Manga").
    Where("manga_id = ?", mangaId).
    Order("language").
    Scan(ctx)
  return util.CheckSliceResult(result, err).Unwrap()
}

func (t translationRepository) FindMangaSpecific(mangaId string, language common.Language) (*mangas.Translation, error) {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  result := new(mangas.Translation)
  err := t.db.NewSelect().
    Model(result).
    Relation("Manga").
    Where("manga_id = ? AND language = ?", mangaId, language).
    Scan(ctx)

  if err != nil {
    return nil, err
  }

  return result, nil
}

func (t translationRepository) FindById(id string) (*mangas.Translation, error) {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  var result = new(mangas.Translation)
  err := t.db.NewSelect().
    Model(result).
    Relation("Manga").
    Where("id = ?", id).
    Order("language").
    Scan(ctx)
  if err != nil {
    return nil, err
  }
  return result, nil
}

func (t translationRepository) Update(translation *mangas.Translation) error {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  res, err := t.db.NewUpdate().
    Model(translation).
    WherePK().
    ExcludeColumn("id", "manga_id").
    Exec(ctx)

  return util.CheckSqlResult(res, err)
}

func (t translationRepository) DeleteByMangaId(mangaId string) error {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  res, err := t.db.NewDelete().
    Model(util.Nil[mangas.Translation]()).
    Where("manga_id = ?", mangaId).
    Exec(ctx)
  return util.CheckSqlResult(res, err)
}

func (t translationRepository) DeleteMangaSpecific(mangaId string, languages []common.Language) error {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  res, err := t.db.NewDelete().
    Model(util.Nil[mangas.Translation]()).
    Where("manga_id = ? AND language IN (?)", mangaId, bun.In(languages)).
    Exec(ctx)

  return util.CheckSqlResult(res, err)
}

func (t translationRepository) DeleteByIds(translationIds []string) error {
  ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
  defer cancel()

  res, err := t.db.NewDelete().
    Model(util.Nil[mangas.Translation]()).
    Where("id IN (?)", bun.In(translationIds)).
    Exec(ctx)
  return util.CheckSqlResult(res, err)
}
