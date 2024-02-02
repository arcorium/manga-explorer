package pg

import (
	"context"
	"github.com/uptrace/bun"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/repository"
	repo "manga-explorer/internal/infrastructure/repository"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/containers"
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
		WherePK().
		ExcludeColumn("id", "translator_id", "created_at").
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
		ColumnExpr("COUNT(DISTINCT comment.*) AS total_comment, chapter.*").
		Join("LEFT JOIN comments AS comment").
		JoinOn("comment.object_type = ?", mangas.CommentObjectChapter.String()).
		JoinOn("comment.object_id = chapter.id").
		Relation("Pages", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Order("page.number")
		}).
		Relation("Translator").
		Where("chapter.id = ?", id).
		Group("chapter.id", "translator.id").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return chapter, nil
}

func (c chapterRepository) FindVolumeDetails(volumeId string) (*mangas.Volume, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	volume := new(mangas.Volume)
	err := c.db.NewSelect().
		Model(volume).
		Relation("Chapters", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Order("chapter.number", "chapter.created_at").
				ColumnExpr("COUNT(DISTINCT comment.*) AS total_comment, chapter.*").
				Join("LEFT JOIN comments AS comment").
				JoinOn("comment.object_type = ?", mangas.CommentObjectChapter.String()).
				JoinOn("comment.object_id = chapter.id").
				Group("chapter.id", "translator.id")
		}).
		Relation("Chapters.Translator").
		Where("volume.id = ?", volumeId).
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return volume, nil
}

func (c chapterRepository) FindPagesDetails(chapterId string, pages []uint16) ([]mangas.Page, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result := []mangas.Page{}
	query := c.db.NewSelect().
		Model(&result).
		Where("chapter_id = ? AND number IN (?)", chapterId, bun.In(pages))

	err := query.Scan(ctx)

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

func (c chapterRepository) InsertChapterHistories(history *mangas.ChapterHistory) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewInsert().
		Model(history).
		On("CONFLICT (user_id, chapter_id) DO UPDATE SET last_view = EXCLUDED.last_view").
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (c chapterRepository) FindMangaChapterHistories(userId string, mangaId string, pagedQuery repo.QueryParameter) (repo.PagedQueryResult[[]mangas.Chapter], error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.ChapterHistory

	query := c.db.NewSelect().
		//Model(util.Nil[mangas.ChapterHistory]()).
		Model(&result).
		ExcludeColumn("*").
		Relation("Chapter").
		Relation("Chapter.Translator").
		Relation("Chapter.Volume").
		Where("user_id = ? AND chapter__volume.manga_id = ?", userId, mangaId).
		OrderExpr("last_view DESC")

	query = pagedQuery.Insert(query)
	count, err := query.ScanAndCount(ctx)

	actuals := containers.CastSlicePtr(result, func(current *mangas.ChapterHistory) mangas.Chapter {
		return *current.Chapter
	})

	res := util.CheckSliceResult(actuals, err)
	return repo.NewResult(res.Data, count), res.Err
}
