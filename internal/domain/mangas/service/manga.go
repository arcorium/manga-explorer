package service

import (
	"manga-explorer/internal/common"
	dto2 "manga-explorer/internal/common/dto"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/mangas/dto"
)

type IManga interface {
	// CreateManga create new manga
	CreateManga(input *dto.MangaCreateInput) status.Object
	UpdateMangaCover(input *dto.MangaCoverUpdateInput) status.Object
	EditManga(input *dto.MangaEditInput) status.Object
	// CreateVolume Upsert a new volume which should be belonged to manga with 0 chapters
	CreateVolume(input *dto.VolumeCreateInput) status.Object
	// DeleteVolume Delete volume and make the chapters based on the volume into NULL
	DeleteVolume(input *dto.VolumeDeleteInput) status.Object
	// CreateComments Upsert a new comment for manga, chapter, and page
	CreateComments(input *dto.MangaCommentCreateInput) status.Object
	// UpsertMangaRating Upsert or Update manga rating
	UpsertMangaRating(input *dto.RateUpsertInput) status.Object
	FindMangaHistories(userId string, query *dto2.PagedQueryInput) ([]dto.MangaHistoryResponse, *dto2.ResponsePage, status.Object)
	FindMangaFavorites(userId string, query *dto2.PagedQueryInput) ([]dto.MangaFavoriteResponse, *dto2.ResponsePage, status.Object)
	ListMangas(query *dto2.PagedQueryInput) ([]dto.MangaResponse, *dto2.ResponsePage, status.Object)
	SearchMangas(query *dto.MangaSearchQuery) ([]dto.MangaResponse, *dto2.ResponsePage, status.Object)
	InsertMangaTranslations(input *dto.MangaInsertTranslationInput) status.Object
	FindMangaTranslations(mangaId string) ([]dto.TranslationResponse, status.Object)
	FindSpecificMangaTranslation(mangaId string, language common.Language) (dto.TranslationResponse, status.Object)
	DeleteMangaTranslations(mangaId string) status.Object
	DeleteTranslations(input *dto.TranslationDeleteInput) status.Object
	UpdateTranslation(input *dto.TranslationUpdateInput) status.Object
	// FindMangaByIds find mangas based on the ids
	FindMangaByIds(mangaId ...string) ([]dto.MangaResponse, status.Object)
	// FindRandomMangas find random based mangas and will return n manga count. n is limit parameter
	FindRandomMangas(limit uint64) ([]dto.MangaResponse, status.Object)
	// FindMangaComments find all manga comments
	FindMangaComments(mangaId string) ([]dto.CommentResponse, status.Object)
	// FindMangaRatings find all manga ratings
	FindMangaRatings(mangaId string) ([]dto.RateResponse, status.Object)
}
