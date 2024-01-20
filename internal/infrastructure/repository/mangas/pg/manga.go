package pg

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/repository"
	repo "manga-explorer/internal/infrastructure/repository"
	"manga-explorer/internal/util"
	"strings"
	"time"
)

func NewManga(db bun.IDB) repository.IManga {
	return &mangaRepository{db: db}
}

type mangaRepository struct {
	db bun.IDB
}

func (m mangaRepository) excludeAllColumns(query *bun.SelectQuery) *bun.SelectQuery {
	return query.ExcludeColumn("*")
}

func (m mangaRepository) mangaRelations(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Relation("Comments").
		Relation("Ratings").
		Relation("Translations").
		Relation("Volumes")
}

func (m mangaRepository) DeleteVolume(mangaId string, volume uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewDelete().
		Model((*mangas.Volume)(nil)).
		Where("manga_id = ? AND number = ?", mangaId, volume).
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (m mangaRepository) CreateManga(mangas *mangas.Manga) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewInsert().
		Model(mangas).
		Returning("NULL").
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (m mangaRepository) EditManga(manga *mangas.Manga) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewUpdate().
		Model(manga).
		WherePK().
		OmitZero().
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (m mangaRepository) FindMangasByFilter(filter *mangas.SearchFilter, pagedQuery repo.QueryParameter) (repo.PagedQueryResult[[]mangas.Manga], error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	//ctx := context.Background()

	var result []mangas.Manga

	var query *bun.SelectQuery
	if filter.HasGenre() {
		query = m.db.NewSelect().
			Model(util.Nil[mangas.MangaGenre]()).
			ExcludeColumn("*").
			Relation("Manga").
			Relation("Genre", m.excludeAllColumns).
			Group("manga.id")

	} else {
		query = m.mangaRelations(m.db.NewSelect().
			Model(util.Nil[mangas.Manga]()))
	}

	if len(filter.Title) != 0 {
		title := "%" + strings.ToLower(filter.Title) + "%"
		query = query.Where("LOWER(manga.original_title) LIKE ?", title)
	}

	if filter.HasOrigin() {
		if filter.IsOriginInclude {
			query = query.Where("manga.origin IN (?)", bun.In(filter.Origins))
		} else {
			query = query.Where("manga.origin NOT IN (?)", bun.In(filter.Origins))
		}
	}

	// Set exclude
	if filter.Genres.HasExclude() {
		exceptQuery := util.Clone(query)
		exceptQuery = exceptQuery.Where("genre.name IN (?)", bun.In(filter.Genres.Exclude))

		query = query.Except(exceptQuery)
	}

	// Has include
	if filter.Genres.HasInclude() {
		query = query.Where("genre.name IN (?)", bun.In(filter.Genres.Include))
		if filter.Genres.IsAndOperation {
			// AND Operation
			query = query.Having("COUNT(genre.name) = ?", len(filter.Genres.Include))
		}
	}

	// Paged
	// TODO: Wrapped query only used for include and exclude genres (has both), for the others perhaps can be removed for maybe gain more performance?
	subquery := m.db.NewSelect().TableExpr("(?)", query) // Wrap it on another query to be able to have LIMIT and OFFSET on outer query
	subquery = pagedQuery.Insert(subquery)

	// Separate the query which Bun also does this on ScanAndCount under the hood
	// Count from main query
	count, err := query.Count(ctx)
	if err != nil {
		return repo.NewResult[[]mangas.Manga](nil, 0), err
	}
	if count == 0 {
		return repo.NewResult[[]mangas.Manga](nil, 0), sql.ErrNoRows
	}

	// Scan from wrapped query
	err = subquery.Scan(ctx, &result)
	if err != nil {
		return repo.NewResult[[]mangas.Manga](nil, 0), err
	}

	if len(result) <= 0 {
		return repo.NewResult[[]mangas.Manga](nil, 0), sql.ErrNoRows
	}

	return repo.NewResult(result, count), nil
}

func (m mangaRepository) FindRandomMangas(limit uint64) ([]mangas.Manga, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.Manga
	query := m.mangaRelations(m.db.NewSelect().
		Model(&result)).
		OrderExpr("RANDOM()").
		Limit(int(limit))

	err := query.Scan(ctx)

	return util.CheckSliceResult(result, err).Unwrap()
}

func (m mangaRepository) ListMangas(parameter repo.QueryParameter) (repo.PagedQueryResult[[]mangas.Manga], error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.Manga
	query := m.mangaRelations(m.db.NewSelect().
		Model(&result))
	query = parameter.Insert(query)

	count, err := query.ScanAndCount(ctx)

	return repo.NewResult(result, uint64(count)), util.CheckSliceResult(result, err).Err
}

func (m mangaRepository) CreateVolume(volume *mangas.Volume) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewInsert().
		Model(volume).
		Returning("NULL").
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (m mangaRepository) FindMangaById(id string) (*mangas.Manga, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result mangas.Manga
	err := m.mangaRelations(m.db.NewSelect().
		Model(&result)).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m mangaRepository) FindMangasById(ids ...string) ([]mangas.Manga, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.Manga
	query := m.mangaRelations(m.db.NewSelect().
		Model(&result)).
		Where("id IN (?)", bun.In(ids))

	err := query.Scan(ctx)

	return util.CheckSliceResult(result, err).Unwrap()
}

func (m mangaRepository) FindMangaHistories(userId string, pagedQuery repo.QueryParameter) (repo.PagedQueryResult[[]mangas.MangaHistory], error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//ctx := context.Background()

	var result []mangas.MangaHistory

	qq := m.db.NewSelect().
		Model(util.Nil[mangas.ChapterHistory]()).
		ColumnExpr("MAX(last_view) as last_view, manga.*").
		Join("JOIN ? AS chapter ON ? = ?", bun.Ident("chapters"), bun.Ident("chapter.id"), bun.Ident("chapter_id")).
		Join("JOIN ? AS volume ON ? = ?", bun.Ident("volumes"), bun.Ident("volume.id"), bun.Ident("chapter.volume_id")).
		Join("JOIN ? AS manga ON ? = ?", bun.Ident("mangas"), bun.Ident("manga.id"), bun.Ident("volume.manga_id")).
		Where("user_id = ?", userId).
		Group("manga.id").
		Order("last_view DESC")

	qq = pagedQuery.Insert(qq)

	count, err := qq.ScanAndCount(ctx, &result)

	res := util.CheckSliceResult(result, err)
	return repo.NewResult(res.Data, count), res.Err
}

func (m mangaRepository) FindMangaFavorites(userId string, pagedQuery repo.QueryParameter) (repo.PagedQueryResult[[]mangas.MangaFavorite], error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//ctx := context.Background()

	var result []mangas.MangaFavorite

	query := m.db.NewSelect().
		Model(&result).
		Where("user_id = ?", userId).
		Relation("Manga").
		Order("manga.original_title")

	query = pagedQuery.Insert(query)
	count, err := query.ScanAndCount(ctx)

	res := util.CheckSliceResult(result, err)
	return repo.NewResult(res.Data, count), res.Err
}

func (m mangaRepository) FindMangaChapterHistories(userId string, mangaId string, pagedQuery repo.QueryParameter) (repo.PagedQueryResult[[]mangas.ChapterHistory], error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.ChapterHistory

	query := m.db.NewSelect().
		Model(&result).
		Where("user_id = ? AND chapter.manga_id = ?", userId, mangaId).
		OrderExpr("last_view DESC").
		Relation("Chapter")

	query = pagedQuery.Insert(query)
	count, err := query.ScanAndCount(ctx)

	res := util.CheckSliceResult(result, err)
	return repo.NewResult(res.Data, count), res.Err
}
