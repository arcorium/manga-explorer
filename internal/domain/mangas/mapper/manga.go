package mapper

import (
	"github.com/biter777/countries"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/infrastructure/file"
	"manga-explorer/internal/util/containers"
	"time"
)

func ToMangaResponse(manga *mangas.Manga) dto.MangaResponse {
	return dto.MangaResponse{
		Id:              manga.Id,
		Status:          manga.Status.Underlying(),
		Origin:          manga.Origin,
		PublicationYear: manga.PublicationYear,
		CoverURL:        manga.CoverURL.HostnameFullpath(file.MangaAsset),
		Comments:        containers.CastSlicePtr(manga.Comments, ToCommentResponse),
		Ratings:         containers.CastSlicePtr(manga.Ratings, ToRatingResponse),
		Translations:    containers.CastSlicePtr(manga.Translations, ToTranslationResponse),
		Volumes:         containers.CastSlicePtr(manga.Volumes, ToVolumeResponse),
		ViewedCount:     0, // TODO: Implement it
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
	status, err := mangas.NewStatus(input.Status)
	manga := mangas.NewManga(input.Title, input.Description, "", input.PublicationYear,
		status, countries.ByName(string(input.Origin)))
	return manga, err
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
		CoverURL:            file.Name(input.CoverUrl),
		UpdatedAt:           time.Now(),
	}, err
}
