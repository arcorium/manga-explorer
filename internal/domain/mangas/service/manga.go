package service

import (
	"manga-explorer/internal/app/common/status"
	appDto "manga-explorer/internal/app/dto"
	"manga-explorer/internal/domain/mangas/dto"
)

type IManga interface {
	// CreateManga create new manga
	CreateManga(input *dto.MangaCreateInput) status.Object
	EditManga(input *dto.MangaEditInput) status.Object
	// CreateVolume Create a new volume which should be belonged to manga with 0 chapters
	CreateVolume(input *dto.VolumeCreateInput) status.Object
	// DeleteVolume Delete volume and make the chapters based on the volume into NULL
	DeleteVolume(input *dto.VolumeDeleteInput) status.Object
	// CreateComments Create a new comment for manga, chapter, and page
	CreateComments(input *dto.MangaCommentCreateInput) status.Object
	// UpsertMangaRating Create or Update manga rating
	UpsertMangaRating(input *dto.RateUpsertInput) status.Object
	// FindPagedMangaHistories works like FindMangaHistories, but it will return in the length of less or equal to limit parameter
	// second return value handle which page is it currently in which can become information for the client
	FindMangaHistories(userId string, query *appDto.PagedQueryInput) ([]dto.MangaHistoryResponse, *appDto.ResponsePage, status.Object)
	// FindPagedMangaFavorites works like FindMangaFavorites, but it will return in the length of less or equal to limit parameter,
	// second return value handle which page is it currently in which can become information for the client
	FindMangaFavorites(userId string, query *appDto.PagedQueryInput) ([]dto.MangaFavoriteResponse, *appDto.ResponsePage, status.Object)
	// ListPagedMangas works like ListMangas, but instead of returning all it will return based on the request
	ListMangas(query *appDto.PagedQueryInput) ([]dto.MangaResponse, *appDto.ResponsePage, status.Object)
	// SearchPagedMangas works like SearchMangas, but it will return n mangas. n is the element count
	SearchPagedMangas(query *dto.MangaSearchQuery) ([]dto.MangaResponse, *appDto.ResponsePage, status.Object)
	// FindMangaByIds find mangas based on the ids
	FindMangaByIds(mangaId ...string) ([]dto.MangaResponse, status.Object)
	// FindRandomMangas find random based mangas and will return n manga count. n is limit parameter
	FindRandomMangas(limit uint64) ([]dto.MangaResponse, status.Object)
	// FindMangaComments find all manga comments
	FindMangaComments(mangaId string) ([]dto.CommentResponse, status.Object)
	// FindMangaRatings find all manga ratings
	FindMangaRatings(mangaId string) ([]dto.RateResponse, status.Object)
}
