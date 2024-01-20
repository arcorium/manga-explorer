package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/util"
	"testing"
)

func createCommentForTest(objectType mangas.CommentObject, objectId, userId, comment string, like, dislike uint64) *mangas.Comment {
	tmp := mangas.NewComment(objectType, objectId, userId, comment, like, dislike)
	return &tmp
}
func createReplyCommentForTest(parentId string, objectType mangas.CommentObject, objectId, userId, comment string, like, dislike uint64) *mangas.Comment {
	tmp := mangas.NewReplyComment(parentId, objectType, objectId, userId, comment, like, dislike)
	return &tmp
}

func Test_commentRepository_CreateComment(t *testing.T) {
	type args struct {
		comment *mangas.Comment
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				comment: createCommentForTest(mangas.CommentObjectPage, "fd9eb653-5158-4e10-bccb-bcfd2b12fd25", "f8b6a114-8cfe-4e14-b2fb-590c53cec0f1", "Hwo", 15, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Self reply",
			args: args{
				comment: createReplyCommentForTest(PageParent.Id, PageParent.ObjectType, PageParent.ObjectId, PageParent.UserId, "Hwo", 15, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Multiple comment single user",
			args: args{
				comment: createCommentForTest(MangaParent.ObjectType, MangaParent.ObjectId, MangaParent.UserId, "Hello", 10, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Bad object id as uuid",
			args: args{
				comment: createCommentForTest(mangas.CommentObjectPage, "asdadszczc-asdas", "f8b6a114-8cfe-4e14-b2fb-590c53cec0f1", "Hwo", 15, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "User not found",
			args: args{
				comment: createCommentForTest(mangas.CommentObjectPage, "fd9eb653-5158-4e10-bccb-bcfd2b12fd25", uuid.NewString(), "Hwo", 15, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
		{
			name: "Bad user id as uuid",
			args: args{
				comment: createCommentForTest(mangas.CommentObjectPage, "fd9eb653-5158-4e10-bccb-bcfd2b12fd25", "asdadaz-zxczc", "Hwo", 15, 1),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)
		c := NewComment(tx)

		t.Run(tt.name, func(t *testing.T) {
			defer func(tx bun.Tx) {
				require.NoError(t, tx.Rollback())
			}(tx)

			err := c.CreateComment(tt.args.comment)
			if !tt.wantErr(t, err) {
				t.Errorf("CreateComment(%v)", tt.args.comment)
				return
			}

			if err != nil {
				return
			}
			// Find created comment
			got, err := c.FindComment(tt.args.comment.Id)
			require.NoError(t, err)
			require.NotNil(t, got)

			// Ignore time
			got.CreatedAt = tt.args.comment.CreatedAt
			got.UpdatedAt = tt.args.comment.UpdatedAt
			// Ignore relations
			got.User = tt.args.comment.User
			got.ParentComment = tt.args.comment.ParentComment

			assert.Equalf(t, tt.args.comment, got, "CreateComment(%v)", tt.args.comment.Id)
		})
	}
}

func Test_commentRepository_DeleteComment(t *testing.T) {
	type args struct {
		commentId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				commentId: MangaParent.Id,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Comment doesn't exists",
			args: args{
				commentId: uuid.NewString(),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Bad comment id as uuid",
			args: args{
				commentId: "asdadsa0zxcxz-asd1231",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		tx, err := Db.Begin()
		require.NoError(t, err)
		c := NewComment(tx)

		t.Run(tt.name, func(t *testing.T) {
			err := c.DeleteComment(tt.args.commentId)
			if !tt.wantErr(t, err) {
				t.Errorf("DeleteComment(%v)", tt.args.commentId)
				return
			}

			if err != nil {
				return
			}
			// Find created comment
			got, err := c.FindComment(tt.args.commentId)
			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}

func Test_commentRepository_FindChapterComments(t *testing.T) {
	type args struct {
		mangaId string
	}
	tests := []struct {
		name    string
		args    args
		want    []mangas.Comment
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				mangaId: ChapterParent.ObjectId,
			},
			want: ChapterComments,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Another kind of object id",
			args: args{
				mangaId: MangaParent.ObjectId,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Chapter object doesn't exists",
			args: args{
				mangaId: uuid.NewString(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Bad chapter id as uuid",
			args: args{
				mangaId: "asdasd12-xzc",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		c := NewComment(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.FindChapterComments(tt.args.mangaId)
			if !tt.wantErr(t, err, fmt.Sprintf("FindChapterComments(%v)", tt.args.mangaId)) {
				return
			}

			if count := util.NilCount[mangas.Comment](got, tt.want); count == 2 {
				return
			} else if count == 1 {
				t.Errorf("FindChapterComments(): expected: %v got: %v", tt.want, got)
				return
			}

			require.Len(t, got, len(tt.want))
			for i := 0; i < len(got); i++ {
				// Ignore time
				got[i].UpdatedAt = tt.want[i].UpdatedAt
				got[i].CreatedAt = tt.want[i].CreatedAt
				// Ignore relations
				got[i].ParentComment = tt.want[i].ParentComment
				got[i].User = tt.want[i].User
			}

			assert.Equalf(t, tt.want, got, "FindChapterComments(%v)", tt.args.mangaId)
		})
	}
}

func Test_commentRepository_FindMangaComments(t *testing.T) {
	type args struct {
		mangaId string
	}
	tests := []struct {
		name    string
		args    args
		want    []mangas.Comment
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				mangaId: MangaParent.ObjectId,
			},
			want: MangaComments,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Another kind of object id",
			args: args{
				mangaId: PageParent.ObjectId,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Manga object doesn't exists",
			args: args{
				mangaId: uuid.NewString(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Bad manga id as uuid",
			args: args{
				mangaId: "asdasd12-xzc",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		c := NewComment(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.FindMangaComments(tt.args.mangaId)
			if !tt.wantErr(t, err) {
				t.Errorf("FindMangaComments(%v)", tt.args.mangaId)
				return
			}

			if count := util.NilCount[mangas.Comment](got, tt.want); count == 2 {
				return
			} else if count == 1 {
				t.Errorf("FindMangaComments(): expected: %v got: %v", tt.want, got)
				return
			}

			require.Len(t, got, len(tt.want))
			for i := 0; i < len(got); i++ {
				// Ignore time
				got[i].UpdatedAt = tt.want[i].UpdatedAt
				got[i].CreatedAt = tt.want[i].CreatedAt
				// Ignore relations
				got[i].ParentComment = tt.want[i].ParentComment
				got[i].User = tt.want[i].User
			}

			assert.Equalf(t, tt.want, got, "FindMangaComments(%v)", tt.args.mangaId)
		})
	}
}

func Test_commentRepository_FindPageComments(t *testing.T) {
	type args struct {
		pageId string
	}
	tests := []struct {
		name    string
		args    args
		want    []mangas.Comment
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				pageId: PageParent.ObjectId,
			},
			want: PageComments,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Another kind of object id",
			args: args{
				pageId: ChapterParent.ObjectId,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Page object doesn't exists",
			args: args{
				pageId: uuid.NewString(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Bad page id as uuid",
			args: args{
				pageId: "asdasd12-xzc",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		c := NewComment(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.FindPageComments(tt.args.pageId)
			if !tt.wantErr(t, err) {
				t.Errorf("FindPageComments(%v)", tt.args.pageId)
				return
			}

			if count := util.NilCount[mangas.Comment](got, tt.want); count == 2 {
				return
			} else if count == 1 {
				t.Errorf("FindPageComments(): expected: %v got: %v", tt.want, got)
				return
			}

			require.Len(t, got, len(tt.want))
			for i := 0; i < len(got); i++ {
				// Ignore time
				got[i].UpdatedAt = tt.want[i].UpdatedAt
				got[i].CreatedAt = tt.want[i].CreatedAt
				// Ignore relations
				got[i].ParentComment = tt.want[i].ParentComment
				got[i].User = tt.want[i].User
			}

			assert.Equalf(t, tt.want, got, "FindPageComments(%v)", tt.args.pageId)
		})
	}
}

func Test_commentRepository_FindComment(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *mangas.Comment
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Normal",
			args: args{
				id: PageParent.Id,
			},
			want: &PageParent,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "Comment doesn't exists",
			args: args{
				id: uuid.NewString(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, sql.ErrNoRows)
			},
		},
		{
			name: "Bad comment id as uuid",
			args: args{
				id: "asdasdb-123bdas",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		c := NewComment(Db)
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.FindComment(tt.args.id)
			if !tt.wantErr(t, err) {
				t.Errorf("FindComment(%v)", tt.args.id)
				return
			}

			if count := util.NilCount[mangas.Comment](got, tt.want); count == 2 {
				return
			} else if count == 1 {
				t.Errorf("FindComment(): expected: %v got: %v", tt.want, got)
				return
			}

			// Ignore time
			got.CreatedAt = tt.want.CreatedAt
			got.UpdatedAt = tt.want.UpdatedAt
			// Ignore relations
			got.User = tt.want.User
			got.ParentComment = tt.want.ParentComment

			assert.Equalf(t, tt.want, got, "FindComment(%v)", tt.args.id)
		})
	}
}
