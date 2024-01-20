package mangas

import (
	"github.com/uptrace/bun"
	"manga-explorer/internal/domain/users"
	"time"
)

type MangaFavorite struct {
	bun.BaseModel `bun:"table:manga_favorites"`

	// Composite primary key
	UserId    string    `bun:",type:uuid,pk"`
	MangaId   string    `bun:",type:uuid,pk"`
	CreatedAt time.Time `bun:",notnull"`

	Manga *Manga      `bun:"rel:belongs-to,join:manga_id=id,on_delete:CASCADE"`
	User  *users.User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
}

func NewFavorite(userId, mangaId string) MangaFavorite {
	currentTime := time.Now()
	return MangaFavorite{
		UserId:    userId,
		MangaId:   mangaId,
		CreatedAt: currentTime,
	}
}
