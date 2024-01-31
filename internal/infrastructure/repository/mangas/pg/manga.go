package pg

import (
	"context"
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

func (m mangaRepository) DeleteVolume(mangaId string, volume uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewDelete().
		Model((*mangas.Volume)(nil)).
		Where("manga_id = ? AND number = ?", mangaId, volume).
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (m mangaRepository) CreateManga(mangas *mangas.Manga, genres []mangas.MangaGenre) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Insert manga
	res, err := tx.NewInsert().
		Model(mangas).
		Returning("NULL").
		Exec(ctx)

	if err != nil {
		tx.Rollback()
		return util.CheckSqlResult(res, err)
	}

	// Insert manga genre
	res, err = tx.NewInsert().
		Model(&genres).
		Returning("NULL").
		Exec(ctx)

	if err != nil {
		tx.Rollback()
		return util.CheckSqlResult(res, err)
	}
	err = tx.Commit()
	return nil
}

func (m mangaRepository) EditManga(manga *mangas.Manga) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewUpdate().
		Model(manga).
		WherePK().
		ExcludeColumn("created_at", "id", "cover_url").
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (m mangaRepository) EditMangaGenres(additionals, removes []mangas.MangaGenre) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	//ctx := context.Background()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Add genres
	if len(additionals) > 0 {
		res, err := tx.NewInsert().
			Model(&additionals).
			Exec(ctx)

		if err != nil {
			tx.Rollback()
			return util.CheckSqlResult(res, err)
		}
	}

	// Remove genres
	if len(removes) > 0 {
		res, err := tx.NewDelete().
			Model(&removes).
			WherePK().
			Exec(ctx)

		if err != nil {
			tx.Rollback()
			return util.CheckSqlResult(res, err)
		}
	}

	return tx.Commit()
}

func (m mangaRepository) FindMangasByFilter(filter *mangas.SearchFilter, pagedQuery repo.QueryParameter) (repo.PagedQueryResult[[]mangas.Manga], error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var result []mangas.Manga

	query := m.getMangaSelectQuery(&result).
		Join("LEFT JOIN manga_genres ON manga_genres.manga_id = manga.id").
		Join("LEFT JOIN genres ON genres.id = manga_genres.genre_id").
		Relation("Genres").
		Order("manga.original_title")

	if len(filter.Title) > 0 {
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

	// Has include
	if filter.Genres.HasInclude() {
		query = query.Where("genres.name IN (?)", bun.In(filter.Genres.Include))
		if filter.Genres.IsAndOperation {
			// AND Operation
			query = query.Having("COUNT(DISTINCT genres.*) >= ?", len(filter.Genres.Include))
		}
	}

	// Set exclude
	if filter.Genres.HasExclude() {
		//exceptQuery := util.Clone(query)
		//exceptQuery = exceptQuery.Where("genres.name IN (?)", bun.In(filter.Genres.Exclude))
		// Prevent using EXCEPT, because bun can't use LIMIT and OFFSET for the result
		exceptQuery := m.db.NewSelect().
			Table("mangas").
			Join("JOIN manga_genres ON manga_genres.manga_id = mangas.id").
			Join("JOIN genres ON genres.id = manga_genres.genre_id").
			ColumnExpr("mangas.id AS id").
			Where("genres.name IN (?)", bun.In(filter.Genres.Exclude)).
			Group("mangas.id")

		query = query.With("result", exceptQuery).
			Where("manga.id NOT IN (SELECT result.id FROM result)")
		//query = query.Where("manga.id NOT IN ?", exceptQuery)
		//query = query.Except(exceptQuery)
	}

	// Paged
	query = pagedQuery.Insert(query)
	count, err := query.ScanAndCount(ctx)

	res := util.CheckSliceResult(result, err)
	return repo.NewResult(res.Data, count), res.Err
}

func (m mangaRepository) getMangaSelectQuery(model any) *bun.SelectQuery {
	return m.db.NewSelect().
		Model(model).
		Join("LEFT JOIN ? ON ? = ?", bun.Ident("rates"), bun.Ident("manga.id"), bun.Ident("manga_id")).
		Join("LEFT JOIN comments AS comment").
		JoinOn("comment.object_type = ?", mangas.CommentObjectManga.String()).
		JoinOn("comment.object_id = manga.id").
		ColumnExpr("AVG(rates.rate) AS average_rate, COUNT(DISTINCT rates.*) AS total_rater").
		ColumnExpr("COUNT(DISTINCT comment.*) AS total_comment").
		ColumnExpr("manga.*").
		Group("manga.id")
}

func (m mangaRepository) FindRandomMangas(limit uint64) ([]mangas.Manga, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.Manga
	query := m.getMangaSelectQuery(&result).
		Relation("Genres").
		OrderExpr("RANDOM()").
		Limit(int(limit)).
		Order("manga.original_title")

	err := query.Scan(ctx)

	return util.CheckSliceResult(result, err).Unwrap()
}

func (m mangaRepository) ListMangas(parameter repo.QueryParameter) (repo.PagedQueryResult[[]mangas.Manga], error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.Manga
	query := m.getMangaSelectQuery(&result).
		Relation("Genres").
		Order("manga.updated_at DESC")
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

func (m mangaRepository) FindMinimalMangaById(id string) (*mangas.Manga, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result mangas.Manga
	query := m.getMangaSelectQuery(&result).
		Group("manga.id").
		Where("id = ", id)

	err := query.Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m mangaRepository) FindMangasById(ids ...string) ([]mangas.Manga, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.Manga
	query := m.getMangaSelectQuery(&result).
		Relation("Genres").
		Relation("Volumes", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Order("number")
		}).
		Relation("Volumes.Chapters", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Order("number")
		}).
		Relation("Volumes.Chapters.Translator").
		Relation("Translations").
		Where("manga.id IN (?)", bun.In(ids)).
		Group("manga.id")

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

func (m mangaRepository) InsertMangaFavorite(favorite *mangas.MangaFavorite) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewInsert().
		Model(favorite).
		Exec(ctx)

	return util.CheckSqlResult(res, err)
}

func (m mangaRepository) RemoveMangaFavorite(favorite *mangas.MangaFavorite) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := m.db.NewDelete().
		Model(favorite).
		WherePK().
		Exec(ctx)

	return util.CheckSqlResult(res, err)
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
