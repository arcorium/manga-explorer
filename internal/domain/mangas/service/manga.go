package service

import (
	"manga-explorer/internal/app/common"
	commonDto "manga-explorer/internal/app/dto"
	"manga-explorer/internal/domain/mangas/dto"
)

type IManga interface {
	// CreateManga create new manga
	CreateManga(input *dto.MangaCreateInput) common.Status
	EditManga(input *dto.MangaEditInput) common.Status
	// CreateVolume Create a new volume which should be belonged to manga with 0 chapters
	CreateVolume(input *dto.VolumeCreateInput) common.Status
	// DeleteVolume Delete volume and make the chapters based on the volume into NULL
	DeleteVolume(input *dto.VolumeDeleteInput) common.Status
	// CreateComments Create a new comment for manga, chapter, and page
	CreateComments(input *dto.MangaCommentCreateInput) common.Status
	// UpsertMangaRating Create or Update manga rating
	UpsertMangaRating(input *dto.RateUpsertInput) common.Status
	FindMangaHistories(userId string) ([]dto.MangaHistoryResponse, common.Status)
	// FindPagedMangaHistories works like FindMangaHistories, but it will return in the length of less or equal to limit parameter
	// second return value handle which page is it currently in which can become information for the client
	FindPagedMangaHistories(userId string, query *commonDto.PagedQueryInput) ([]dto.MangaHistoryResponse, commonDto.ResponsePage, common.Status)
	// FindMangaFavorites find all user favorites manga
	FindMangaFavorites(userId string) ([]dto.MangaFavoriteResponse, common.Status)
	// FindPagedMangaFavorites works like FindMangaFavorites, but it will return in the length of less or equal to limit parameter,
	// second return value handle which page is it currently in which can become information for the client
	FindPagedMangaFavorites(userId string, query *commonDto.PagedQueryInput) ([]dto.MangaFavoriteResponse, commonDto.ResponsePage, common.Status)
	// ListMangas get all manga (it can be slow when the data is huge)
	ListMangas() ([]dto.MangaResponse, common.Status)
	// ListPagedMangas works like ListMangas, but instead of returning all it will return based on the request
	ListPagedMangas(query *commonDto.PagedQueryInput) ([]dto.MangaResponse, commonDto.ResponsePage, common.Status)
	// SearchMangas find manga based on the title and the genres, both is nullable
	SearchMangas(query *dto.MangaSearchQuery) ([]dto.MangaResponse, common.Status)
	// SearchPagedMangas works like SearchMangas, but it will return n mangas. n is the element count
	SearchPagedMangas(query *dto.MangaSearchQuery) ([]dto.MangaResponse, commonDto.ResponsePage, common.Status)
	// FindMangaByIds find mangas based on the ids
	FindMangaByIds(mangaId ...string) ([]dto.MangaResponse, common.Status)
	// FindRandomMangas find random based mangas and will return n manga count. n is limit parameter
	FindRandomMangas(limit uint64) ([]dto.MangaResponse, common.Status)
	// FindMangaComments find all manga comments
	FindMangaComments(mangaId string) ([]dto.CommentResponse, common.Status)
	// FindMangaRatings find all manga ratings
	FindMangaRatings(mangaId string) ([]dto.RateResponse, common.Status)
}
