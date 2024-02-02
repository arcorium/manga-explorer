package dto

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/common"
	"manga-explorer/internal/common/dto"
	"mime/multipart"
	"time"
)

type MangaResponse struct {
	Id              string         `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"desc"`
	Status          string         `json:"status"`
	Origin          common.Country `json:"origin"`
	PublicationYear uint16         `json:"year"`
	CoverURL        string         `json:"cover_url"`

	Rate         float32 `json:"rate"`
	TotalRater   uint64  `json:"total_rater"`
	TotalComment uint64  `json:"total_comment"`
	//Comments     []CommentResponse     `json:"comments,omitempty"`
	//Ratings      []RateResponse        `json:"ratings,omitempty"`
	Translations []TranslationResponse `json:"translations,omitempty"`
	Volumes      []VolumeResponse      `json:"volumes,omitempty"`
	Genres       []GenreResponse       `json:"genres"`
}

type MinimalMangaResponse struct {
	Id              string          `json:"id"`
	Title           string          `json:"title"`
	Description     string          `json:"desc"`
	Status          string          `json:"status"`
	Origin          common.Country  `json:"origin"`
	PublicationYear uint16          `json:"year"`
	CoverURL        string          `json:"cover_url"`
	Rate            float32         `json:"rate"`
	TotalRater      uint64          `json:"total_rater"`
	TotalComment    uint64          `json:"total_comment"`
	Genres          []GenreResponse `json:"genres"`
}

type MangaHistoryResponse struct {
	MangaResponse
	LastView time.Time `json:"last_view"`
}

type MangaFavoriteResponse struct {
	MangaResponse
	FavoritedAt time.Time `json:"favorited_at"`
}

type MangaCreateInput struct {
	Title           string         `json:"title" binding:"required"`
	Description     string         `json:"desc" binding:"required"`
	Status          string         `json:"status" binding:"required,manga_status"`
	Origin          common.Country `json:"origin" binding:"required,iso3166_1_alpha3|iso3166_1_alpha2"`
	PublicationYear uint16         `json:"publication_year" binding:"required"`
	Genres          []string       `json:"genres" binding:"required,dive,uuid4"`
}

type MangaCoverUpdateInput struct {
	MangaId string                `uri:"manga_id" binding:"required,uuid4"`
	Image   *multipart.FileHeader `form:"image" binding:"required"`
}

func (c *MangaCoverUpdateInput) ConstructURI(ctx *gin.Context) {
	c.MangaId = ctx.Param("manga_id")
}

type MangaEditInput struct {
	MangaId         string         `uri:"manga_id" binding:"required,uuid4"`
	Status          string         `json:"status" binding:"required,manga_status"`
	Origin          common.Country `json:"origin" binding:"required,iso3166_1_alpha3|iso3166_1_alpha2"`
	Title           string         `json:"title" binding:"required,min=1"`
	Description     string         `json:"description"`
	PublicationYear uint16         `json:"publication_year" binding:"required"`
}

type MangaGenreEditInput struct {
	MangaId       string   `uri:"manga_id" binding:"required,uuid4"`
	AddGenres     []string `json:"adds" binding:"omitempty,dive,uuid4"`
	RemovedGenres []string `json:"removes" binding:"omitempty,dive,uuid4"`
}

func (m *MangaGenreEditInput) ConstructURI(ctx *gin.Context) {
	m.MangaId = ctx.Param("manga_id")
}

func (e *MangaEditInput) ConstructURI(ctx *gin.Context) {
	e.MangaId = ctx.Param("manga_id")
}

type MangaSearchQuery struct {
	dto.PagedQueryInput
	Title  string                              `json:"title"`
	Genres common.CriterionOption[string]      `json:"genre"`
	Origin common.IncludeArray[common.Country] `json:"origin"`
}

type FavoriteMangaInput struct {
	Operator string `json:"op" binding:"required,oneof=add remove"`
	UserId   string `json:"-"`
	MangaId  string `json:"manga_id" binding:"required,uuid4"`
}

type MangaChapterHistoriesFindInput struct {
	UserId  string `json:"-"`
	MangaId string `uri:"manga_id" binding:"required,uuid4"`
	dto.PagedQueryInput
}

func (c *MangaChapterHistoriesFindInput) ConstructURI(ctx *gin.Context) {
	c.MangaId = ctx.Param("manga_id")
}
