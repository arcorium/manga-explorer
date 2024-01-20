package mapper

import (
	"github.com/google/uuid"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/util/containers"
	"time"
)

func ToMangaResponse(manga *mangas.Manga) dto.MangaResponse {
	return dto.MangaResponse{
		Status:          manga.Status.Underlying(),
		Origin:          manga.Origin,
		PublicationYear: manga.PublicationYear,
		Comments:        containers.CastSlicePtr(manga.Comments, ToCommentResponse),
		Ratings:         containers.CastSlicePtr(manga.Ratings, ToRatingResponse),
		Translations:    containers.CastSlicePtr(manga.Translations, ToTranslationResponse),
		Volumes:         containers.CastSlicePtr(manga.Volumes, ToVolumeResponse),
		ViewedCount:     0,
		FavoriteCount:   0,
	}
}

func ToMangaHistoryResponse(history *mangas.MangaHistory) dto.MangaHistoryResponse {
	return dto.MangaHistoryResponse{
		MangaResponse: ToMangaResponse(history.Manga),
		LastView:      history.LastView,
	}
}

func ToMangaFavoriteResponse(favorite *mangas.MangaFavorite) dto.MangaFavoriteResponse {
	return dto.MangaFavoriteResponse{
		MangaResponse: ToMangaResponse(favorite.Manga),
		FavoritedAt:   favorite.CreatedAt,
	}
}

func MapMangaCreateInput(input *dto.MangaCreateInput) (mangas.Manga, error) {
	now := time.Now()
	status, err := mangas.NewStatus(input.Status)
	return mangas.Manga{
		Id:              uuid.NewString(),
		Status:          status,
		Origin:          input.Origin,
		PublicationYear: input.PublicationYear,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, err
}

func MapMangaEditInput(input *dto.MangaEditInput) (mangas.Manga, error) {
	status, err := mangas.NewStatus(input.Status)
	return mangas.Manga{
		Id:                  input.MangaId,
		Status:              status,
		Origin:              input.Origin,
		OriginalTitle:       input.Title,
		OriginalDescription: input.Description,
		PublicationYear:     input.PublicationYear,
		CoverURL:            input.CoverUrl,
		UpdatedAt:           time.Now(),
	}, err
}
