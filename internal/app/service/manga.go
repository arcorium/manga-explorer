package service

import (
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/app/dto"
	appMapper "manga-explorer/internal/app/mapper"
	"manga-explorer/internal/domain/mangas"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
	"manga-explorer/internal/infrastructure/file"
	fileService "manga-explorer/internal/infrastructure/file/service"
	"manga-explorer/internal/util/containers"
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
	return status.ConditionalRepository(err, status.CREATED)
}

func (m mangaService) DeleteVolume(input *mangaDto.VolumeDeleteInput) status.Object {
	err := m.mangaRepo.DeleteVolume(input.MangaId, input.Volume)
	return status.ConditionalRepository(err, status.DELETED)
}

func (m mangaService) CreateComments(input *mangaDto.MangaCommentCreateInput) status.Object {
	comment := mapper.MapMangaCommentCreateInput(input)
	if input.HasParent() {
		parent, err := m.commentRepo.FindComment(input.ParentId)
		if err != nil {
			return status.InternalError()
		}

		// Validate
		if !comment.ValidateAsReply(parent) {
			return status.Error(status.BAD_BODY_REQUEST_ERROR)
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return status.ConditionalRepository(err, status.CREATED)
}

func (m mangaService) UpsertMangaRating(input *mangaDto.RateUpsertInput) status.Object {
	rate := mapper.MapRateUpsertInput(input)
	err := m.rateRepo.Upsert(&rate)
	return status.ConditionalRepository(err, status.SUCCESS)
}

func (m mangaService) ListMangas(query *dto.PagedQueryInput) ([]mangaDto.MangaResponse, *dto.ResponsePage, status.Object) {
	result, err := m.mangaRepo.ListMangas(query.ToQueryParam())
	if err != nil {
		return nil, nil, status.RepositoryError(err)
	}
	mangaResponses := containers.CastSlicePtr1(result.Data, m.fileService, mapper.ToMangaResponse)
	responsePage := appMapper.NewResponsePage(mangaResponses, result.Total, query)
	return mangaResponses, &responsePage, status.Success()
}

func (m mangaService) SearchMangas(query *mangaDto.MangaSearchQuery) ([]mangaDto.MangaResponse, *dto.ResponsePage, status.Object) {
	filter := mapper.MapMangaSearchQuery(query)
	res, err := m.mangaRepo.FindMangasByFilter(&filter, query.ToQueryParam())
	if err != nil {
		return nil, nil, status.RepositoryError(err)
	}
	mangaResponses := containers.CastSlicePtr1(res.Data, m.fileService, mapper.ToMangaResponse)
	responsePage := appMapper.NewResponsePage(mangaResponses, res.Total, &query.PagedQueryInput)
	return mangaResponses, &responsePage, status.Success()
}

func (m mangaService) FindMangaByIds(mangaIds ...string) ([]mangaDto.MangaResponse, status.Object) {
	manga, err := m.mangaRepo.FindMangasById(mangaIds...)
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	return containers.CastSlicePtr1(manga, m.fileService, mapper.ToMangaResponse), status.Success()
}

func (m mangaService) FindRandomMangas(limit uint64) ([]mangaDto.MangaResponse, status.Object) {
	lisMangas, err := m.mangaRepo.FindRandomMangas(limit)
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	mangaResponses := containers.CastSlicePtr1(lisMangas, m.fileService, mapper.ToMangaResponse)
	return mangaResponses, status.Success()
}

func (m mangaService) FindMangaComments(mangaId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindMangaComments(mangaId)
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	commentResponses := containers.CastSlicePtr(comments, mapper.ToCommentResponse)
	return commentResponses, status.Success()
}

func (m mangaService) FindMangaRatings(mangaId string) ([]mangaDto.RateResponse, status.Object) {
	ratings, err := m.rateRepo.FindMangaRatings(mangaId)
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	ratingResponses := containers.CastSlicePtr(ratings, mapper.ToRatingResponse)
	return ratingResponses, status.Success()
}

func (m mangaService) CreateManga(input *mangaDto.MangaCreateInput) status.Object {
	model, err := mapper.MapMangaCreateInput(input)
	if err != nil {
		return status.Error(status.BAD_BODY_REQUEST_ERROR)
	}

	err = m.mangaRepo.CreateManga(&model)
	return status.ConditionalRepository(err, status.CREATED)
}

func (m mangaService) UpdateMangaCover(input *mangaDto.MangaCoverUpdateInput) status.Object {
	manga, err := m.mangaRepo.FindMangaById(input.MangaId)
	if err != nil {
		return status.RepositoryError(err)
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
	err = m.mangaRepo.EditManga(&editedManga)
	return status.ConditionalRepository(err, status.UPDATED)
}

func (m mangaService) EditManga(input *mangaDto.MangaEditInput) status.Object {
	model, err := mapper.MapMangaEditInput(input)
	if err != nil {
		return status.Error(status.BAD_BODY_REQUEST_ERROR)
	}

	err = m.mangaRepo.EditManga(&model)
	return status.ConditionalRepository(err, status.UPDATED)
}

func (m mangaService) InsertMangaTranslations(input *mangaDto.MangaInsertTranslationInput) status.Object {
	translates := mapper.MapInsertTranslateInput(input)

	err := m.translationRepo.Create(translates)
	return status.ConditionalRepository(err, status.CREATED)
}

func (m mangaService) FindMangaTranslations(mangaId string) ([]mangaDto.TranslationResponse, status.Object) {
	translates, err := m.translationRepo.FindByMangaId(mangaId)
	if err != nil {
		return nil, status.RepositoryError(err)
	}
	return containers.CastSlicePtr(translates, mapper.ToTranslationResponse), status.Success()
}

func (m mangaService) FindSpecificMangaTranslation(mangaId string, language common.Language) (mangaDto.TranslationResponse, status.Object) {
	// Convert into 3 letter Country code
	lang := common.NewLanguage(language.Code())
	translate, err := m.translationRepo.FindMangaSpecific(mangaId, lang)
	if err != nil {
		return mangaDto.TranslationResponse{}, status.RepositoryError(err)
	}
	return mapper.ToTranslationResponse(translate), status.Success()
}

func (m mangaService) DeleteMangaTranslations(mangaId string) status.Object {
	err := m.translationRepo.DeleteByMangaId(mangaId)
	return status.ConditionalRepository(err, status.DELETED)
}

func (m mangaService) DeleteTranslations(input *mangaDto.TranslationDeleteInput) status.Object {
	err := m.translationRepo.DeleteByIds(input.TranslationIds)
	return status.ConditionalRepository(err, status.DELETED)
}

func (m mangaService) UpdateTranslation(input *mangaDto.TranslationUpdateInput) status.Object {
	translate := mapper.MapTranslationUpdateInput(input)

	err := m.translationRepo.Update(&translate)
	return status.ConditionalRepository(err, status.UPDATED)
}

func (m mangaService) FindMangaHistories(userId string, query *dto.PagedQueryInput) ([]mangaDto.MangaHistoryResponse, *dto.ResponsePage, status.Object) {
	res, err := m.mangaRepo.FindMangaHistories(userId, query.ToQueryParam())
	if err != nil {
		return nil, nil, status.RepositoryError(err)
	}

	mangaResult := containers.CastSlicePtr1(res.Data, m.fileService, mapper.ToMangaHistoryResponse)
	pages := appMapper.NewResponsePage(mangaResult, res.Total, query)
	return mangaResult, &pages, status.Success()
}

func (m mangaService) FindMangaFavorites(userId string, query *dto.PagedQueryInput) ([]mangaDto.MangaFavoriteResponse, *dto.ResponsePage, status.Object) {
	res, err := m.mangaRepo.FindMangaFavorites(userId, query.ToQueryParam())
	if err != nil {
		return nil, nil, status.RepositoryError(err)
	}

	mangaResponses := containers.CastSlicePtr1(res.Data, m.fileService, mapper.ToMangaFavoriteResponse)
	pages := appMapper.NewResponsePage(mangaResponses, res.Total, query)
	return mangaResponses, &pages, status.Success()
}
