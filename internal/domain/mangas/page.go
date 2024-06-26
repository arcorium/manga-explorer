package mangas

import (
  "github.com/google/uuid"
  "manga-explorer/internal/infrastructure/file"
)

type Page struct {
  Id        string    `bun:",pk,type:uuid"`
  ChapterId string    `bun:",nullzero,notnull,unique:page_chapter_idx,type:uuid"`
  Number    uint16    `bun:",nullzero,notnull,unique:page_chapter_idx"`
  ImageURL  file.Name `bun:",notnull,nullzero"`

  Chapter *Chapter `bun:"rel:belongs-to,join:chapter_id=id,on_delete:CASCADE"`
}

func NewPage(chapterId string, filename file.Name, number uint16) Page {
  return Page{
    Id:        uuid.NewString(),
    ChapterId: chapterId,
    Number:    number,
    ImageURL:  filename,
  }
}
