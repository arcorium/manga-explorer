package mapper

import (
	"github.com/google/uuid"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/users/mapper"
	"time"
)

func ToCommentResponse(comment *mangas.Comment) dto.CommentResponse {
	return dto.CommentResponse{
		User:     mapper.ToUserResponse(comment.User),
		Comment:  comment.Comment,
		IsEdited: comment.IsEdited,
		Like:     comment.Like,
		Dislike:  comment.Dislike,
	}
}

func MapMangaCommentCreateInput(input *dto.MangaCommentCreateInput) mangas.Comment {
	now := time.Now()
	return mangas.Comment{
		Id:         uuid.NewString(),
		ParentId:   input.ParentId,
		ObjectType: mangas.CommentObjectManga,
		ObjectId:   input.MangaId,
		UserId:     input.UserId,
		Comment:    input.Comment,
		IsEdited:   false,
		Like:       0,
		Dislike:    0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func MapChapterCommentCreateInput(input *dto.ChapterCommentCreateInput) mangas.Comment {
	now := time.Now()
	return mangas.Comment{
		Id:         uuid.NewString(),
		ParentId:   input.ParentId,
		ObjectType: mangas.CommentObjectChapter,
		ObjectId:   input.ChapterId,
		UserId:     input.UserId,
		Comment:    input.Comment,
		IsEdited:   false,
		Like:       0,
		Dislike:    0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
func MapPageCommentCreateInput(input *dto.PageCommentCreateInput) mangas.Comment {
	now := time.Now()
	return mangas.Comment{
		Id:         uuid.NewString(),
		ParentId:   input.ParentId,
		ObjectType: mangas.CommentObjectPage,
		ObjectId:   input.PageId,
		UserId:     input.UserId,
		Comment:    input.Comment,
		IsEdited:   false,
		Like:       0,
		Dislike:    0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
