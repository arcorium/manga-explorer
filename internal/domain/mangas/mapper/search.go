package mapper

import (
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
)

func MapMangaSearchQuery(query *dto.MangaSearchQuery) mangas.SearchFilter {
	return mangas.SearchFilter{
		Title:           query.Title,
		Genres:          query.Genres,
		Origins:         query.Origin.Values,
		IsOriginInclude: query.Origin.IsInclude,
	}
}
