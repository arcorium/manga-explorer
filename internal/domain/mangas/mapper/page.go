package mapper

import (
	"github.com/google/uuid"
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
)

func ToPageResponse(page *mangas.Page) dto.PageResponse {
	return dto.PageResponse{
		Page:     page.Number,
		ImageURL: page.ImageURL,
	}
}

func MapPageCreateInput(input *dto.PageCreateInput) mangas.Page {
	return mangas.Page{
		Id:        uuid.NewString(),
		ChapterId: input.ChapterId,
		Number:    input.Page,
	}
}
