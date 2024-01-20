package pg

import (
	"context"
	"github.com/uptrace/bun"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/util"
	"time"
)

func NewComment(db bun.IDB) repository.IComment {
	return &commentRepository{db: db}
}

type commentRepository struct {
	db bun.IDB
}

func (c commentRepository) findComments(objectType string, objectId string) ([]mangas.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var result []mangas.Comment
	query := c.db.NewSelect().
		Model(&result).
		Where("object_type = ? AND object_id = ?", objectType, objectId)
	err := query.Scan(ctx)
	return util.CheckSliceResult(result, err).Unwrap()
}

func (c commentRepository) FindMangaComments(mangaId string) ([]mangas.Comment, error) {
	return c.findComments(mangas.CommentObjectManga.String(), mangaId)
}

func (c commentRepository) FindChapterComments(chapterId string) ([]mangas.Comment, error) {
	return c.findComments(mangas.CommentObjectChapter.String(), chapterId)
}

func (c commentRepository) FindPageComments(pageId string) ([]mangas.Comment, error) {
	return c.findComments(mangas.CommentObjectPage.String(), pageId)
}

func (c commentRepository) CreateComment(comment *mangas.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewInsert().
		Model(comment).
		Returning("NULL").
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (c commentRepository) DeleteComment(commentId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.db.NewDelete().
		Model((*mangas.Comment)(nil)).
		Where("id = ?", commentId).
		Exec(ctx)
	return util.CheckSqlResult(res, err)
}

func (c commentRepository) FindComment(id string) (*mangas.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res := new(mangas.Comment)
	err := c.db.NewSelect().
		Model(res).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return res, nil
}
