package mapper

import (
  "github.com/google/uuid"
  "manga-explorer/internal/domain/mangas"
  "manga-explorer/internal/domain/mangas/dto"
  "manga-explorer/internal/domain/users/mapper"
  fileService "manga-explorer/internal/infrastructure/file/service"
  "manga-explorer/internal/util/containers"
  "time"
)

func ToChapterResponse(chapter *mangas.Chapter, fs fileService.IFile) dto.ChapterResponse {
  return dto.ChapterResponse{
    Id:           chapter.Id,
    Language:     chapter.Language.ParseLang(),
    Chapter:      chapter.Number,
    Title:        chapter.Title,
    TotalComment: &chapter.TotalComment,
    CreatedAt:    chapter.CreatedAt,
    Comments:     containers.CastSlicePtr(chapter.Comments, toCommentResponse),
    Pages:        containers.CastSlicePtr1(chapter.Pages, fs, ToPageResponse),
    Translator:   mapper.ToUserResponse(chapter.Translator),
  }
}

func ToMinimalChapterResponse(chapter *mangas.Chapter) dto.ChapterResponse {
  return dto.ChapterResponse{
    Id:         chapter.Id,
    Language:   chapter.Language.ParseLang(),
    Chapter:    chapter.Number,
    Title:      chapter.Title,
    CreatedAt:  chapter.CreatedAt,
    Comments:   containers.CastSlicePtr(chapter.Comments, toCommentResponse),
    Translator: mapper.ToUserResponse(chapter.Translator),
  }
}

func MapChapterCreateInput(input *dto.ChapterCreateInput) mangas.Chapter {
  now := time.Now()
  chapter := mangas.Chapter{
    Id:           uuid.NewString(),
    VolumeId:     input.VolumeId,
    Language:     input.Language.ParseLang(),
    Title:        input.Title,
    TranslatorId: input.TranslatorId,
    PublishDate:  input.PublishDate,
    Number:       input.Chapter,
    CreatedAt:    now,
    UpdatedAt:    now,
  }

  return chapter
}

func MapChapterEditInput(input *dto.ChapterEditInput) mangas.Chapter {
  return mangas.Chapter{
    Id:          input.ChapterId,
    VolumeId:    input.VolumeId,
    Language:    input.Language.ParseLang(),
    Title:       input.Title,
    Number:      input.Chapter,
    PublishDate: input.PublishDate,
    UpdatedAt:   time.Now(),
  }
}
