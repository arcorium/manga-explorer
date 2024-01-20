package mapper

import (
	"github.com/google/uuid"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/users/mapper"
	"manga-explorer/internal/util/containers"
	"time"
)

func ToChapterResponse(chapter *mangas.Chapter) dto.ChapterResponse {
	return dto.ChapterResponse{
		Language:   common.Country(chapter.Language),
		Title:      chapter.Title,
		CreatedAt:  chapter.CreatedAt,
		Comments:   containers.CastSlicePtr(chapter.Comments, ToCommentResponse),
		Pages:      containers.CastSlicePtr(chapter.Pages, ToPageResponse),
		Translator: mapper.ToUserResponse(chapter.Translator),
	}
}

func MapChapterCreateInput(input *dto.ChapterCreateInput) mangas.Chapter {
	now := time.Now()
	chapter := mangas.Chapter{
		Id:           uuid.NewString(),
		VolumeId:     input.VolumeId,
		Language:     common.Language(input.Language),
		Title:        input.Title,
		TranslatorId: input.TranslatorId,
		PublishDate:  input.PublishDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return chapter
}

func MapChapterEditInput(input *dto.ChapterEditInput) mangas.Chapter {
	return mangas.Chapter{
		Id:          input.ChapterId,
		VolumeId:    input.VolumeId,
		Language:    common.Language(input.Language),
		Title:       input.Title,
		PublishDate: input.PublishDate,
		UpdatedAt:   time.Now(),
	}
}
