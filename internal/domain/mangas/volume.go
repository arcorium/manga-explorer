package mangas

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Volume struct {
	bun.BaseModel `bun:"table:volumes"`

	Id          string `bun:",type:uuid,pk"`
	MangaId     string `bun:",unique:manga_volume_idx,type:uuid"`
	Number      uint32 `bun:",notnull,unique:manga_volume_idx"`
	Title       string `bun:",nullzero,"`
	Description string `bun:",nullzero,type:text"`
	//CoverImageURL string `bun:",nullzero"`

	Manga    *Manga    `bun:"rel:belongs-to,join:manga_id=id,on_delete:CASCADE"`
	Chapters []Chapter `bun:"rel:has-many,join:id=volume_id"`
}

func NewVolume(mangaId string, number uint32, title, desc string) Volume {
	return Volume{
		Id:          uuid.NewString(),
		MangaId:     mangaId,
		Number:      number,
		Title:       title,
		Description: desc,
	}
}
