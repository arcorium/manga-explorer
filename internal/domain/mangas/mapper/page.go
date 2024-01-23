package mapper

import (
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
	"manga-explorer/internal/infrastructure/file"
)

func ToPageResponse(page *mangas.Page) dto.PageResponse {
	return dto.PageResponse{
		Page:     page.Number,
		ImageURL: page.ImageURL.HostnameFullpath(file.MangaAsset),
	}
}

func MapPageCreateInput(input *dto.PageCreateInput, filenames []file.Name) []mangas.Page {
	if len(input.Pages) != len(filenames) {
		return nil
	}
	pages := make([]mangas.Page, 0, len(filenames))
	for i := 0; i < len(input.Pages); i++ {
		pages = append(pages, mangas.NewPage(input.ChapterId, filenames[i], input.Pages[i].Number))
	}
	return pages
}
