package mapper

import (
  "manga-explorer/internal/domain/mangas"
  "manga-explorer/internal/domain/mangas/dto"
  "manga-explorer/internal/infrastructure/file"
  fileService "manga-explorer/internal/infrastructure/file/service"
)

func ToPageResponse(page *mangas.Page, fs fileService.IFile) dto.PageResponse {
  return dto.PageResponse{
    Page:     page.Number,
    ImageURL: fs.GetFullpath(file.MangaAsset, page.ImageURL),
  }
}

//func MapPageCreateInput(input *dto.PageCreateInput, filenames []file.Name) []mangas.Page {
//	if len(input.Page) != len(filenames) {
//		return nil
//	}
//	pages := make([]mangas.Page, 0, len(filenames))
//	for i := 0; i < len(input.Page); i++ {
//		pages = append(pages, mangas.NewPage(input.ChapterId, filenames[i], input.Page[i].Number))
//	}
//	return pages
//}

func MapPageCreateInput(input *dto.PageCreateInput, filename file.Name) mangas.Page {
  return mangas.NewPage(input.ChapterId, filename, input.Page.Number)
}
