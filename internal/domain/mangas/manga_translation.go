package mangas

import (
	"github.com/biter777/countries"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"manga-explorer/internal/app/common"
)

type Translation struct {
	bun.BaseModel `bun:"table:manga_translations"`

	Id          string          `bun:",pk,type:uuid"`
	MangaId     string          `bun:",nullzero,notnull,type:uuid"`
	Language    common.Language `bun:",notnull,type:varchar(3)"`
	Title       string          `bun:",nullzero,unique,notnull"`
	Description string          `bun:",type:text"`

	Manga *Manga `bun:"rel:belongs-to,join:manga_id=id,on_delete:CASCADE"`
}

func NewTranslation(mangaId, title, desc string, lang countries.CountryCode) Translation {
	return Translation{
		Id:          uuid.NewString(),
		MangaId:     mangaId,
		Language:    common.Language(lang.Alpha3()),
		Title:       title,
		Description: desc,
	}
}
