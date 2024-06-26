package mapper

import (
  "manga-explorer/internal/domain/mangas"
  "manga-explorer/internal/domain/mangas/dto"
)

func ToTranslationResponse(translation *mangas.Translation) dto.TranslationResponse {
  return dto.TranslationResponse{
    Id:          translation.Id,
    Lang:        translation.Language,
    Title:       translation.Title,
    Description: translation.Description,
  }
}

func MapCreateTranslateInput(input *dto.InternalTranslation, mangaId string) mangas.Translation {
  return mangas.NewTranslation(mangaId, input.Title, input.Description, input.Lang.ParseLang())
}

func MapInsertTranslateInput(input *dto.MangaTranslationInsertInput) []mangas.Translation {
  result := make([]mangas.Translation, 0, len(input.Translations))

  for i := 0; i < len(input.Translations); i++ {
    result = append(result, MapCreateTranslateInput(&input.Translations[i], input.MangaId))
  }
  return result
}

func MapTranslationUpdateInput(input *dto.TranslationEditInput) mangas.Translation {
  return mangas.Translation{
    Id:          input.TranslationId,
    Language:    input.Lang.ParseLang(),
    Title:       input.Title,
    Description: input.Description,
  }
}
