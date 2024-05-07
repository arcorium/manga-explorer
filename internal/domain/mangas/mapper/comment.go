package mapper

import (
  "github.com/google/uuid"
  "manga-explorer/internal/domain/mangas"
  "manga-explorer/internal/domain/mangas/dto"
  "manga-explorer/internal/domain/users/mapper"
  "time"
)

func toCommentResponse(comment *mangas.Comment) dto.CommentResponse {
  return dto.CommentResponse{
    Id:       comment.Id,
    User:     mapper.ToUserResponse(comment.User),
    Comment:  comment.Comment,
    IsEdited: comment.IsEdited,
    Like:     comment.Like,
    Dislike:  comment.Dislike,
    Time:     comment.CreatedAt,
  }
}

func ToCommentsResponse(comments []mangas.Comment) []dto.CommentResponse {
  type Wrapper struct {
    Index int
  }

  var result []dto.CommentResponse
  var ids = map[string][]Wrapper{} // index on replies

  for i := 0; i < len(comments); i++ {
    val := &comments[i]

    // Base comment
    if len(val.ParentId) == 0 {
      result = append(result, toCommentResponse(val))
      ids[val.Id] = []Wrapper{{i}}
    } else {
      // Replies
      parentIndices := ids[val.ParentId]
      var parent = &result[parentIndices[0].Index]
      for j := 1; j < len(parentIndices); j++ {
        parent = &parent.Replies[parentIndices[j].Index]
      }

      // Set as child
      parent.Replies = append(parent.Replies, toCommentResponse(val))
      insertedIndex := len(parent.Replies) - 1

      // Add to indexes
      copied := make([]Wrapper, len(parentIndices))
      copy(copied, parentIndices)
      copied = append(copied, Wrapper{insertedIndex})
      ids[val.Id] = copied
    }
  }
  return result
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
