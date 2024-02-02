package service

import (
	"manga-explorer/internal/common"
	commonDto "manga-explorer/internal/common/dto"
	appMapper "manga-explorer/internal/common/mapper"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/domain/mangas"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/infrastructure/file"
	fileService "manga-explorer/internal/infrastructure/file/service"
	"manga-explorer/internal/util/containers"
	"manga-explorer/internal/util/opt"
	"time"
)

func NewMangaService(fileService fileService.IFile, mangaRepo repository.IManga, translation repository.ITranslation, commentRepo repository.IComment, rateRepo repository.IRate) service.IManga {
	return &mangaService{
		fileService:     fileService,
		mangaRepo:       mangaRepo,
		commentRepo:     commentRepo,
		rateRepo:        rateRepo,
		translationRepo: translation,
	}
}

type mangaService struct {
	fileService fileService.IFile

	mangaRepo       repository.IManga
	translationRepo repository.ITranslation
	commentRepo     repository.IComment
	rateRepo        repository.IRate
}

func (m mangaService) CreateVolume(input *mangaDto.VolumeCreateInput) status.Object {
	volume := mapper.MapVolumeCreateInput(input)

	err := m.mangaRepo.CreateVolume(&volume)
	return status.ConditionalRepositoryE(err, status.CREATED, opt.New(status.VOLUME_ALREADY_EXISTS), opt.New(status.VOLUME_CREATE_FAILED))
}

func (m mangaService) DeleteVolume(input *mangaDto.VolumeDeleteInput) status.Object {
	err := m.mangaRepo.DeleteVolume(input.MangaId, input.Volume)
	return status.ConditionalRepository(err, status.DELETED, opt.New(status.VOLUME_DELETE_FAILED))
}

func (m mangaService) CreateComments(input *mangaDto.MangaCommentCreateInput) status.Object {
	comment := mapper.MapMangaCommentCreateInput(input)
	if input.HasParent() {
		parent, err := m.commentRepo.FindComment(input.ParentId)
		if err != nil {
			return status.RepositoryError(err, opt.New(status.COMMENT_PARENT_NOT_FOUND))
		}

		// Response
		if !comment.ValidateAsReply(parent) {
			return status.Error(status.COMMENT_PARENT_DIFFERENT_SCOPE)
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return status.ConditionalRepository(err, status.CREATED, opt.New(status.COMMENT_CREATE_FAILED))
}

func (m mangaService) UpsertMangaRating(input *mangaDto.RateUpsertInput) status.Object {
	rate := mapper.MapRateUpsertInput(input)
	err := m.rateRepo.Upsert(&rate)
	return status.ConditionalRepository(err, status.SUCCESS, opt.New(status.RATING_NOT_FOUND))
}

func (m mangaService) ListMangas(query *commonDto.PagedQueryInput) ([]mangaDto.MinimalMangaResponse, *commonDto.ResponsePage, status.Object) {
	result, err := m.mangaRepo.ListMangas(query.ToQueryParam())
	mangaResponses := containers.CastSlicePtr1(result.Data, m.fileService, mapper.ToMinimalMangaResponse)
	responsePage := appMapper.NewResponsePage(mangaResponses, result.Total, query)
	return mangaResponses, &responsePage, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaService) SearchMangas(query *mangaDto.MangaSearchQuery) ([]mangaDto.MinimalMangaResponse, *commonDto.ResponsePage, status.Object) {
	filter := mapper.MapMangaSearchQuery(query)
	res, err := m.mangaRepo.FindMangasByFilter(&filter, query.ToQueryParam())
	mangaResponses := containers.CastSlicePtr1(res.Data, m.fileService, mapper.ToMinimalMangaResponse)
	responsePage := appMapper.NewResponsePage(mangaResponses, res.Total, &query.PagedQueryInput)
	return mangaResponses, &responsePage, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaService) FindMangaByIds(mangaIds ...string) ([]mangaDto.MangaResponse, status.Object) {
	manga, err := m.mangaRepo.FindMangasById(mangaIds...)
	responses := containers.CastSlicePtr1(manga, m.fileService, mapper.ToMangaResponse)
	return responses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaService) FindRandomMangas(limit uint64) ([]mangaDto.MinimalMangaResponse, status.Object) {
	mangaList, err := m.mangaRepo.FindRandomMangas(limit)
	responses := containers.CastSlicePtr1(mangaList, m.fileService, mapper.ToMinimalMangaResponse)
	return responses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaService) FindMangaComments(mangaId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindMangaComments(mangaId)
	responses := mapper.ToCommentsResponse(comments)
	return responses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaService) FindMangaRatings(mangaId string) ([]mangaDto.RateResponse, status.Object) {
	ratings, err := m.rateRepo.FindMangaRatings(mangaId)
	responses := containers.CastSlicePtr(ratings, mapper.ToRatingResponse)
	return responses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaService) CreateManga(input *mangaDto.MangaCreateInput) status.Object {
	model, genres, err := mapper.MapMangaCreateInput(input)
	if err != nil {
		return status.Error(status.BAD_REQUEST_ERROR)
	}

	err = m.mangaRepo.CreateManga(&model, genres)
	return status.ConditionalRepository(err, status.CREATED, opt.New(status.MANGA_CREATE_ALREADY_EXIST))
}

func (m mangaService) UpdateMangaCover(input *mangaDto.MangaCoverUpdateInput) status.Object {
	manga, err := m.mangaRepo.FindMinimalMangaById(input.MangaId)
	if err != nil {
		return status.RepositoryError(err, opt.New(status.MANGA_NOT_FOUND))
	}

	// Upload new cover image
	filename, stat := m.fileService.Upload(file.CoverAsset, input.Image)
	if stat.IsError() {
		return stat
	}

	// Delete current cover image
	if len(manga.CoverURL) != 0 {
		stat = m.fileService.Delete(file.CoverAsset, manga.CoverURL)
		if stat.IsError() {
			return stat
		}
	}

	// Update metadata
	editedManga := mangas.Manga{Id: manga.Id, CoverURL: filename, UpdatedAt: time.Now()}
	err = m.mangaRepo.PatchManga(&editedManga)
	return status.ConditionalRepository(err, status.UPDATED, opt.New(status.MANGA_UPDATE_FAILED))
}

func (m mangaService) EditManga(input *mangaDto.MangaEditInput) status.Object {
	model, err := mapper.MapMangaEditInput(input)
	if err != nil {
		return status.Error(status.BAD_REQUEST_ERROR)
	}

	err = m.mangaRepo.EditManga(&model)
	return status.ConditionalRepository(err, status.UPDATED, opt.New(status.MANGA_UPDATE_FAILED))
}

func (m mangaService) EditMangaGenres(input *mangaDto.MangaGenreEditInput) status.Object {
	additionals, removes := mapper.MapMangaGenreEditInput(input)

	err := m.mangaRepo.EditMangaGenres(additionals, removes)
	return status.ConditionalRepositoryE(err, status.UPDATED, opt.New(status.MANGA_UPDATE_FAILED), opt.New(status.MANGA_UPDATE_FAILED))
}

func (m mangaService) InsertMangaTranslations(input *mangaDto.MangaInsertTranslationInput) status.Object {
	translates := mapper.MapInsertTranslateInput(input)

	err := m.translationRepo.Create(translates)
	return status.ConditionalRepositoryE(err, status.CREATED, opt.New(status.MANGA_TRANSLATION_ALREADY_EXIST), opt.New(status.MANGA_TRANSLATION_CREATE_FAILED))
}

func (m mangaService) FindMangaTranslations(mangaId string) ([]mangaDto.TranslationResponse, status.Object) {
	translates, err := m.translationRepo.FindByMangaId(mangaId)
	responses := containers.CastSlicePtr(translates, mapper.ToTranslationResponse)
	return responses, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaService) FindSpecificMangaTranslation(mangaId string, language common.Language) (mangaDto.TranslationResponse, status.Object) {
	translate, err := m.translationRepo.FindMangaSpecific(mangaId, language)
	if err != nil {
		return mangaDto.TranslationResponse{}, status.RepositoryError(err, opt.New(status.SUCCESS))
	}
	response := mapper.ToTranslationResponse(translate)
	return response, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))

}

func (m mangaService) DeleteMangaTranslations(input *mangaDto.TranslationMangaDeleteInput) status.Object {
	err := m.translationRepo.DeleteMangaSpecific(input.MangaId, input.Languages)
	return status.ConditionalRepository(err, status.DELETED, opt.New(status.MANGA_HAS_NO_TRANSLATIONS))
}

func (m mangaService) DeleteTranslations(input *mangaDto.TranslationDeleteInput) status.Object {
	err := m.translationRepo.DeleteByIds(input.TranslationIds)
	return status.ConditionalRepository(err, status.DELETED, opt.New(status.MANGA_TRANSLATION_NOT_FOUND))
}

func (m mangaService) UpdateTranslation(input *mangaDto.TranslationUpdateInput) status.Object {
	translate := mapper.MapTranslationUpdateInput(input)

	err := m.translationRepo.Update(&translate)
	return status.ConditionalRepository(err, status.UPDATED, opt.New(status.MANGA_TRANSLATION_UPDATE_FAILED))
}

func (m mangaService) FindMangaHistories(userId string, query *commonDto.PagedQueryInput) ([]mangaDto.MangaHistoryResponse, *commonDto.ResponsePage, status.Object) {
	res, err := m.mangaRepo.FindMangaHistories(userId, query.ToQueryParam())

	responses := containers.CastSlicePtr1(res.Data, m.fileService, mapper.ToMangaHistoryResponse)
	pages := appMapper.NewResponsePage(responses, res.Total, query)
	return responses, &pages, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaService) FindMangaFavorites(userId string, query *commonDto.PagedQueryInput) ([]mangaDto.MangaFavoriteResponse, *commonDto.ResponsePage, status.Object) {
	res, err := m.mangaRepo.FindMangaFavorites(userId, query.ToQueryParam())

	responses := containers.CastSlicePtr1(res.Data, m.fileService, mapper.ToMangaFavoriteResponse)
	pages := appMapper.NewResponsePage(responses, res.Total, query)
	return responses, &pages, status.ConditionalRepository(err, status.SUCCESS, opt.New(status.SUCCESS))
}

func (m mangaService) AddFavoriteManga(input *mangaDto.FavoriteMangaInput) status.Object {
	model := mapper.MapFavoriteMangaInput(input)

	err := m.mangaRepo.InsertMangaFavorite(&model)
	return status.ConditionalRepositoryE(err, status.CREATED, opt.New(status.MANGA_UPDATE_FAILED), opt.New(status.MANGA_UPDATE_FAILED))
}

func (m mangaService) RemoveFavoriteManga(input *mangaDto.FavoriteMangaInput) status.Object {
	model := mapper.MapFavoriteMangaInput(input)

	err := m.mangaRepo.RemoveMangaFavorite(&model)
	return status.ConditionalRepositoryE(err, status.CREATED, opt.New(status.MANGA_UPDATE_FAILED), opt.New(status.MANGA_UPDATE_FAILED))
}
