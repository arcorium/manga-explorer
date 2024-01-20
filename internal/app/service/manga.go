package service

import (
	"log"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/app/dto"
	appMapper "manga-explorer/internal/app/mapper"
	"manga-explorer/internal/app/service/utility/file"
	mangaDto "manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/domain/mangas/mapper"
	"manga-explorer/internal/domain/mangas/repository"
	"manga-explorer/internal/domain/mangas/service"
	repo "manga-explorer/internal/infrastructure/repository"
	"manga-explorer/internal/util/containers"
)

func NewMangaService(fileService file.IService, mangaRepo repository.IManga, commentRepo repository.IComment, rateRepo repository.IRate) service.IManga {
	return &mangaService{
		fileService: fileService,
		mangaRepo:   mangaRepo,
		commentRepo: commentRepo,
		rateRepo:    rateRepo,
	}
}

type mangaService struct {
	fileService file.IService

	mangaRepo   repository.IManga
	commentRepo repository.IComment
	rateRepo    repository.IRate
}

func (m mangaService) CreateVolume(input *mangaDto.VolumeCreateInput) common.Status {
	volume := mapper.MapVolumeCreateInput(input)

	err := m.mangaRepo.CreateVolume(&volume)
	return common.NewRepositoryStatus(err, status.SUCCESS_CREATED)
}

func (m mangaService) DeleteVolume(input *mangaDto.VolumeDeleteInput) common.Status {
	err := m.mangaRepo.DeleteVolume(input.MangaId, input.Volume)
	return common.NewRepositoryStatus(err)
}
func (m mangaService) CreateComments(input *mangaDto.MangaCommentCreateInput) common.Status {
	comment := mapper.MapMangaCommentCreateInput(input)
	if input.HasParent() {
		parent, err := m.commentRepo.FindComment(input.ParentId)
		if err != nil {
			return common.StatusError(status.INTERNAL_SERVER_ERROR)
		}

		// Validate
		if !comment.ValidateAsReply(parent) {
			return common.StatusError(status.BAD_BODY_REQUEST_ERROR)
		}
	}
	err := m.commentRepo.CreateComment(&comment)
	return common.NewRepositoryStatus(err, status.SUCCESS_CREATED)
}

func (m mangaService) UpsertMangaRating(input *mangaDto.RateUpsertInput) common.Status {
	rate := mapper.MapRateUpsertInput(input)
	err := m.rateRepo.Upsert(&rate)
	return common.NewRepositoryStatus(err)
}

func (m mangaService) listMangas(parameter repo.QueryParameter) ([]mangaDto.MangaResponse, uint64, common.Status) {
	result, err := m.mangaRepo.ListMangas(parameter)
	if err != nil {
		return nil, 0, common.NewRepositoryStatus(err)
	}
	mangaResponses := containers.CastSlicePtr(result.Data, mapper.ToMangaResponse)
	return mangaResponses, result.Total, common.StatusSuccess()
}
func (m mangaService) ListMangas() ([]mangaDto.MangaResponse, common.Status) {
	mangaResponses, _, status := m.listMangas(repo.NoQueryParameter)
	return mangaResponses, status
}

func (m mangaService) ListPagedMangas(query *dto.PagedQueryInput) ([]mangaDto.MangaResponse, dto.ResponsePage, common.Status) {
	mangaResponses, total, status := m.listMangas(repo.QueryParameter{Offset: query.Offset(), Limit: query.Element})
	responsePage := appMapper.NewResponsePage(mangaResponses, total, query)
	return mangaResponses, responsePage, status
}

func (m mangaService) searchMangas(query *mangaDto.MangaSearchQuery) ([]mangaDto.MangaResponse, uint64, common.Status) {
	filter := mapper.MapMangaSearchQuery(query)
	res, err := m.mangaRepo.FindMangasByFilter(&filter, repo.QueryParameter{query.Offset(), query.Element})
	if err != nil {
		return nil, 0, common.NewRepositoryStatus(err)
	}
	mangaResponses := containers.CastSlicePtr(res.Data, mapper.ToMangaResponse)
	return mangaResponses, res.Total, common.StatusSuccess()
}
func (m mangaService) SearchMangas(query *mangaDto.MangaSearchQuery) ([]mangaDto.MangaResponse, common.Status) {
	mangaResponses, _, status := m.searchMangas(query)
	return mangaResponses, status
}

func (m mangaService) SearchPagedMangas(query *mangaDto.MangaSearchQuery) ([]mangaDto.MangaResponse, dto.ResponsePage, common.Status) {
	mangaResponses, totalElements, status := m.searchMangas(query)
	responsePage := appMapper.NewResponsePage(mangaResponses, totalElements, &query.PagedQueryInput)
	return mangaResponses, responsePage, status
}

func (m mangaService) FindMangaByIds(mangaIds ...string) ([]mangaDto.MangaResponse, common.Status) {
	manga, err := m.mangaRepo.FindMangasById(mangaIds...)
	if err != nil {
		return nil, common.NewRepositoryStatus(err)
	}
	return containers.CastSlicePtr(manga, mapper.ToMangaResponse), common.StatusSuccess()
}

func (m mangaService) FindRandomMangas(limit uint64) ([]mangaDto.MangaResponse, common.Status) {
	lisMangas, err := m.mangaRepo.FindRandomMangas(limit)
	if err != nil {
		return nil, common.NewRepositoryStatus(err)
	}
	mangaResponses := containers.CastSlicePtr(lisMangas, mapper.ToMangaResponse)
	return mangaResponses, common.StatusSuccess()
}

func (m mangaService) FindMangaComments(mangaId string) ([]mangaDto.CommentResponse, common.Status) {
	comments, err := m.commentRepo.FindMangaComments(mangaId)
	if err != nil {
		return nil, common.NewRepositoryStatus(err)
	}
	commentResponses := containers.CastSlicePtr(comments, mapper.ToCommentResponse)
	return commentResponses, common.StatusSuccess()
}

func (m mangaService) FindMangaRatings(mangaId string) ([]mangaDto.RateResponse, common.Status) {
	ratings, err := m.rateRepo.FindMangaRatings(mangaId)
	if err != nil {
		return nil, common.NewRepositoryStatus(err)
	}
	ratingResponses := containers.CastSlicePtr(ratings, mapper.ToRatingResponse)
	return ratingResponses, common.StatusSuccess()
}

func (m mangaService) CreateManga(input *mangaDto.MangaCreateInput) common.Status {
	model, err := mapper.MapMangaCreateInput(input)
	if err != nil {
		return common.StatusError(status.BAD_BODY_REQUEST_ERROR)
	}

	err = m.mangaRepo.CreateManga(&model)
	return common.NewRepositoryStatus(err, status.SUCCESS_CREATED)
}

func (m mangaService) EditManga(input *mangaDto.MangaEditInput) common.Status {
	model, err := mapper.MapMangaEditInput(input)
	if err != nil {
		return common.StatusError(status.BAD_BODY_REQUEST_ERROR)
	}

	err = m.mangaRepo.EditManga(&model)
	return common.NewRepositoryStatus(err)
}

func (m mangaService) getMangasById(ids ...string) ([]mangaDto.MangaResponse, common.Status) {
	// Get the actual mangaHistories
	mangaResponses, stat := m.FindMangaByIds(ids...)

	if stat.IsError() {
		return nil, stat
	}

	// Check length
	if len(mangaResponses) != len(ids) {
		log.Println("User transaction manga have different length with the actual mangas")
		return nil, common.StatusError(status.INTERNAL_SERVER_ERROR)
	}
	return mangaResponses, stat
}

func (m mangaService) findMangaHistories(userId string, offset, limit uint64) ([]mangaDto.MangaHistoryResponse, uint64, common.Status) {
	// Get manga histories
	res, err := m.mangaRepo.FindMangaHistories(userId, repo.QueryParameter{offset, limit})
	if err != nil {
		return nil, 0, common.NewRepositoryStatus(err)
	}

	mangaResult := containers.CastSlicePtr(res.Data, mapper.ToMangaHistoryResponse)
	if err != nil {
		return nil, 0, common.StatusError(status.INTERNAL_SERVER_ERROR)
	}

	return mangaResult, res.Total, common.StatusSuccess()
}
func (m mangaService) FindMangaHistories(userId string) ([]mangaDto.MangaHistoryResponse, common.Status) {
	historyResponses, _, cerr := m.findMangaHistories(userId, 0, 0)
	return historyResponses, cerr
}
func (m mangaService) FindPagedMangaHistories(userId string, query *dto.PagedQueryInput) ([]mangaDto.MangaHistoryResponse, dto.ResponsePage, common.Status) {
	historyResponses, totalElements, cerr := m.findMangaHistories(userId, query.Offset(), query.Element)
	if cerr.IsError() {
		return nil, dto.ResponsePage{}, cerr
	}
	pages := appMapper.NewResponsePage(historyResponses, totalElements, query)
	return historyResponses, pages, cerr
}
func (m mangaService) findMangaFavorites(userId string, offset, limit uint64) ([]mangaDto.MangaFavoriteResponse, uint64, common.Status) {
	// Get manga favorites
	res, err := m.mangaRepo.FindMangaFavorites(userId, repo.QueryParameter{Offset: offset, Limit: limit})
	cerr := common.NewRepositoryStatus(err)
	if cerr.IsError() {
		return nil, 0, cerr
	}

	mangaResponses := containers.CastSlicePtr(res.Data, mapper.ToMangaFavoriteResponse)
	return mangaResponses, res.Total, cerr
}
func (m mangaService) FindMangaFavorites(userId string) ([]mangaDto.MangaFavoriteResponse, common.Status) {
	mangaFavorites, _, cerr := m.findMangaFavorites(userId, 0, 0)
	return mangaFavorites, cerr
}
func (m mangaService) FindPagedMangaFavorites(userId string, query *dto.PagedQueryInput) ([]mangaDto.MangaFavoriteResponse, dto.ResponsePage, common.Status) {
	historyResponses, totalElements, cerr := m.findMangaFavorites(userId, query.Offset(), query.Element)
	if cerr.IsError() {
		return nil, dto.ResponsePage{}, cerr
	}
	pages := appMapper.NewResponsePage(historyResponses, totalElements, query)
	return historyResponses, pages, cerr
}
