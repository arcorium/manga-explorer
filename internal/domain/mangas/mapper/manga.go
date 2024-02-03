package mapper

import (
	"github.com/biter777/countries"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/infrastructure/file"
	fileService "manga-explorer/internal/infrastructure/file/service"
	"manga-explorer/internal/util/containers"
	"time"
)

func ToMangaResponse(manga *mangas.Manga, fs fileService.IFile) dto.MangaResponse {
	return dto.MangaResponse{
		Id:              manga.Id,
		Title:           manga.OriginalTitle,
		Description:     manga.OriginalDescription,
		Status:          manga.Status.String(),
		Origin:          manga.Origin,
		PublicationYear: manga.PublicationYear,
		CoverURL:        fs.GetFullpath(file.CoverAsset, manga.CoverURL),
		Rate:            manga.AverageRate,
		TotalRater:      manga.TotalRater,
		TotalComment:    manga.TotalComment,
		Translations:    containers.CastSlicePtr(manga.Translations, ToTranslationResponse),
		Volumes:         containers.CastSlicePtr1(manga.Volumes, fs, ToVolumeResponse),
		Genres:          containers.CastSlicePtr(manga.Genres, ToGenreResponse),
	}
}

func ToMinimalMangaResponse(manga *mangas.Manga, iFile fileService.IFile) dto.MinimalMangaResponse {
	return dto.MinimalMangaResponse{
		Id:              manga.Id,
		Title:           manga.OriginalTitle,
		Description:     manga.OriginalDescription,
		Status:          manga.Status.String(),
		Origin:          manga.Origin,
		PublicationYear: manga.PublicationYear,
		CoverURL:        iFile.GetFullpath(file.MangaAsset, manga.CoverURL),
		Rate:            manga.AverageRate,
		TotalRater:      manga.TotalRater,
		TotalComment:    manga.TotalComment,
		Genres:          containers.CastSlicePtr(manga.Genres, ToGenreResponse),
	}
}

func ToMangaHistoryResponse(history *mangas.MangaHistory, fl fileService.IFile) dto.MangaHistoryResponse {
	return dto.MangaHistoryResponse{
		MangaResponse: ToMangaResponse(history.Manga, fl),
		LastView:      history.LastView,
	}
}

func ToMangaFavoriteResponse(favorite *mangas.MangaFavorite, fl fileService.IFile) dto.MangaFavoriteResponse {
	return dto.MangaFavoriteResponse{
		MangaResponse: ToMangaResponse(favorite.Manga, fl),
		FavoritedAt:   favorite.CreatedAt,
	}
}

func MapMangaCreateInput(input *dto.MangaCreateInput) (mangas.Manga, []mangas.MangaGenre, error) {
	status, err := mangas.NewStatus(input.Status)
	manga := mangas.NewManga(input.Title, input.Description, "", input.PublicationYear,
		status, countries.ByName(string(input.Origin)))

	genres := []mangas.MangaGenre{}
	for _, v := range input.Genres {
		genres = append(genres, mangas.NewMangaGenre(manga.Id, v))
	}

	return manga, genres, err
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
		UpdatedAt:           time.Now(),
	}, err
}

func MapMangaGenreEditInput(input *dto.MangaGenreEditInput) (additionals []mangas.MangaGenre, removes []mangas.MangaGenre) {
	for _, v := range input.AddGenres {
		additionals = append(additionals, mangas.NewMangaGenre(input.MangaId, v))
	}

	for _, v := range input.RemovedGenres {
		removes = append(removes, mangas.NewMangaGenre(input.MangaId, v))
	}
	return additionals, removes
}

func MapFavoriteMangaInput(input *dto.FavoriteMangaModificationInput) mangas.MangaFavorite {
	return mangas.NewFavorite(input.UserId, input.MangaId)
}
