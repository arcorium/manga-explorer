package pg

import (
	"context"
	"github.com/uptrace/bun"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/util"
	"time"
)

func NewMangaChapter(db bun.IDB) repository.IChapter {
	return &chapterRepository{db: db}
}

type chapterRepository struct {
	db bun.IDB
}

func (c chapterRepository) CreateChapter(chapter *mangas.Chapter) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewInsert().
		Model(chapter).
		Returning("NULL").
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (c chapterRepository) EditChapter(chapter *mangas.Chapter) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewUpdate().
		Model(chapter).
		OmitZero().
		WherePK().
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (c chapterRepository) DeleteChapter(chapterId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewDelete().
		Model((*mangas.Chapter)(nil)).
		Where("id = ?", chapterId).
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (c chapterRepository) FindChapter(id string) (*mangas.Chapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	chapter := new(mangas.Chapter)
	err := c.db.NewSelect().
		Model(chapter).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return chapter, nil
}

func (c chapterRepository) FindChapterPages(chapterId string) ([]mangas.Page, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.Page
	err := c.db.NewSelect().
		Model(&result).
		Relation("Chapter").
		Where("chapter_id = ?", chapterId).
		Order("number").
		Scan(ctx)
	return util.CheckSliceResult(result, err).Unwrap()
}

func (c chapterRepository) FindVolumeChapters(volumeId string) ([]mangas.Chapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.Chapter
	err := c.db.NewSelect().
		Model(&result).
		Relation("Comments").
		Relation("Pages").
		Relation("Translator").
		Relation("Volume").
		Where("volume_id = ?", volumeId).
		Order("number", "created_at").
		Scan(ctx)
	return util.CheckSliceResult(result, err).Unwrap()
}

func (c chapterRepository) DeleteChapterPages(chapterId string, pages []uint16) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := c.db.NewDelete().
		Model(util.Nil[mangas.Page]()).
		Where("chapter_id = ? AND number IN (?)", chapterId, bun.In(pages))
	res, err := query.Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (c chapterRepository) InsertChapterPages(pages []mangas.Page) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewInsert().
		Model(&pages).
		Returning("NULL").
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}
