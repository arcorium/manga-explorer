package mapper

import (
	"manga-explorer/internal/domain/mangas"
	"manga-explorer/internal/domain/mangas/dto"
)

func ToTranslationResponse(translation *mangas.Translation) dto.TranslationResponse {
	return dto.TranslationResponse{
		Lang:        translation.Language,
		Title:       translation.Title,
		Description: translation.Description,
	}
}
