package mangas

import (
  "github.com/google/uuid"
  "github.com/uptrace/bun"
  "manga-explorer/internal/common"
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

func NewTranslation(mangaId, title, desc string, lang common.Language) Translation {
  return Translation{
    Id:          uuid.NewString(),
    MangaId:     mangaId,
    Language:    lang,
    Title:       title,
    Description: desc,
  }
}
