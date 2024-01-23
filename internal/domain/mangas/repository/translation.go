package repository

import "manga-explorer/internal/domain/mangas"

type ITranslation interface {
	Create(translation []mangas.Translation) error
	FindByMangaId(mangaId string) ([]mangas.Translation, error)
	FindById(id string) (*mangas.Translation, error)
	Update(translation *mangas.Translation) error
	DeleteByMangaId(mangaId string) error
	DeleteByIds(translationIds []string) error
}
