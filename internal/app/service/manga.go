package service

import (
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/app/dto"
	appMapper "manga-explorer/internal/app/mapper"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
	fileService "manga-explorer/internal/infrastructure/file/service"
	"manga-explorer/internal/util/containers"
)

func NewMangaService(fileService fileService.IFile, mangaRepo repository.IManga, commentRepo repository.IComment, rateRepo repository.IRate) service.IManga {
	return &mangaService{
		fileService: fileService,
		mangaRepo:   mangaRepo,
		commentRepo: commentRepo,
		rateRepo:    rateRepo,
	}
}

type mangaService struct {
	fileService fileService.IFile

	mangaRepo   repository.IManga
	commentRepo repository.IComment
	rateRepo    repository.IRate
}

func (m mangaService) CreateVolume(input *mangaDto.VolumeCreateInput) status.Object {
	volume := mapper.MapVolumeCreateInput(input)

	err := m.mangaRepo.CreateVolume(&volume)
	return status.FromRepository(err, status.CREATED)
}

func (m mangaService) DeleteVolume(input *mangaDto.VolumeDeleteInput) status.Object {
	err := m.mangaRepo.DeleteVolume(input.MangaId, input.Volume)
	return status.FromRepository(err, status.DELETED)
}
func (m mangaService) CreateComments(input *mangaDto.MangaCommentCreateInput) status.Object {
	comment := mapper.MapMangaCommentCreateInput(input)
	if input.HasParent() {
		parent, err := m.commentRepo.FindComment(input.ParentId)
		if err != nil {
			return status.Error(status.INTERNAL_SERVER_ERROR)
		}

		// Validate
		if !comment.ValidateAsReply(parent) {
			return status.Error(status.BAD_BODY_REQUEST_ERROR)
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return status.FromRepository(err, status.CREATED)
}

func (m mangaService) UpsertMangaRating(input *mangaDto.RateUpsertInput) status.Object {
	rate := mapper.MapRateUpsertInput(input)
	err := m.rateRepo.Upsert(&rate)
	return status.FromRepository(err)
}

func (m mangaService) ListMangas(query *dto.PagedQueryInput) ([]mangaDto.MangaResponse, *dto.ResponsePage, status.Object) {
	result, err := m.mangaRepo.ListMangas(query.ToQueryParam())
	stat := status.FromRepository(err)
	if stat.IsError() {
		return nil, nil, stat
	}
	mangaResponses := containers.CastSlicePtr(result.Data, mapper.ToMangaResponse)
	responsePage := appMapper.NewResponsePage(mangaResponses, result.Total, query)
	return mangaResponses, &responsePage, stat
}

func (m mangaService) SearchPagedMangas(query *mangaDto.MangaSearchQuery) ([]mangaDto.MangaResponse, *dto.ResponsePage, status.Object) {
	filter := mapper.MapMangaSearchQuery(query)
	res, err := m.mangaRepo.FindMangasByFilter(&filter, query.ToQueryParam())
	stat := status.FromRepository(err)
	if stat.IsError() {
		return nil, nil, stat
	}
	mangaResponses := containers.CastSlicePtr(res.Data, mapper.ToMangaResponse)
	responsePage := appMapper.NewResponsePage(mangaResponses, res.Total, &query.PagedQueryInput)
	return mangaResponses, &responsePage, stat
}

func (m mangaService) FindMangaByIds(mangaIds ...string) ([]mangaDto.MangaResponse, status.Object) {
	manga, err := m.mangaRepo.FindMangasById(mangaIds...)
	stat := status.FromRepository(err)
	if err != nil {
		return nil, stat
	}
	return containers.CastSlicePtr(manga, mapper.ToMangaResponse), stat
}

func (m mangaService) FindRandomMangas(limit uint64) ([]mangaDto.MangaResponse, status.Object) {
	lisMangas, err := m.mangaRepo.FindRandomMangas(limit)
	stat := status.FromRepository(err)
	if err != nil {
		return nil, stat
	}
	mangaResponses := containers.CastSlicePtr(lisMangas, mapper.ToMangaResponse)
	return mangaResponses, stat
}

func (m mangaService) FindMangaComments(mangaId string) ([]mangaDto.CommentResponse, status.Object) {
	comments, err := m.commentRepo.FindMangaComments(mangaId)
	stat := status.FromRepository(err)
	if err != nil {
		return nil, stat
	}
	commentResponses := containers.CastSlicePtr(comments, mapper.ToCommentResponse)
	return commentResponses, stat
}

func (m mangaService) FindMangaRatings(mangaId string) ([]mangaDto.RateResponse, status.Object) {
	ratings, err := m.rateRepo.FindMangaRatings(mangaId)
	stat := status.FromRepository(err)
	if err != nil {
		return nil, stat
	}
	ratingResponses := containers.CastSlicePtr(ratings, mapper.ToRatingResponse)
	return ratingResponses, stat
}

func (m mangaService) CreateManga(input *mangaDto.MangaCreateInput) status.Object {
	model, err := mapper.MapMangaCreateInput(input)
	if err != nil {
		return status.Error(status.BAD_BODY_REQUEST_ERROR)
	}

	err = m.mangaRepo.CreateManga(&model)
	return status.FromRepository(err, status.CREATED)
}

func (m mangaService) EditManga(input *mangaDto.MangaEditInput) status.Object {
	model, err := mapper.MapMangaEditInput(input)
	if err != nil {
		return status.Error(status.BAD_BODY_REQUEST_ERROR)
	}

	err = m.mangaRepo.EditManga(&model)
	return status.FromRepository(err, status.UPDATED)
}

func (m mangaService) FindMangaHistories(userId string, query *dto.PagedQueryInput) ([]mangaDto.MangaHistoryResponse, *dto.ResponsePage, status.Object) {
	res, err := m.mangaRepo.FindMangaHistories(userId, query.ToQueryParam())
	stat := status.FromRepository(err, status.UPDATED)
	if stat.IsError() {
		return nil, nil, stat
	}

	mangaResult := containers.CastSlicePtr(res.Data, mapper.ToMangaHistoryResponse)
	pages := appMapper.NewResponsePage(mangaResult, res.Total, query)
	return mangaResult, &pages, stat
}

func (m mangaService) FindMangaFavorites(userId string, query *dto.PagedQueryInput) ([]mangaDto.MangaFavoriteResponse, *dto.ResponsePage, status.Object) {
	res, err := m.mangaRepo.FindMangaFavorites(userId, query.ToQueryParam())
	cerr := status.FromRepository(err)
	if cerr.IsError() {
		return nil, nil, cerr
	}

	mangaResponses := containers.CastSlicePtr(res.Data, mapper.ToMangaFavoriteResponse)
	pages := appMapper.NewResponsePage(mangaResponses, res.Total, query)
	return mangaResponses, &pages, cerr
}
