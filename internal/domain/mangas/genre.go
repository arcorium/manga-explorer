package mangas

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Genre struct {
	bun.BaseModel `bun:"table:genres"`
	Id            string `bun:",pk,type:uuid"`
	Name          string `bun:",nullzero,notnull,unique"`

	CreatedAt time.Time `bun:",notnull"`
	Mangas    []Manga   `bun:"m2m:manga_genres,join:Genre=Manga"`
}

func NewGenre(name string) Genre {
	return Genre{
		Id:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
	}
}

// MangaGenre used for genres on each manga
type MangaGenre struct {
	bun.BaseModel `bun:"table:manga_genres"`
	MangaId       string `bun:",pk,type:uuid"`
	GenreId       string `bun:",pk,type:uuid"`

	Manga *Manga `bun:"rel:belongs-to,join:manga_id=id,on_delete:CASCADE"`
	Genre *Genre `bun:"rel:belongs-to,join:genre_id=id,on_delete:CASCADE"`
}

func NewMangaGenre(mangaId, genreId string) MangaGenre {
	return MangaGenre{
		MangaId: mangaId,
		GenreId: genreId,
	}
}
