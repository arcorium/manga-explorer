package mangas

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	user_entity "manga-explorer/internal/domain/users"
	"time"
)

type Rate struct {
	bun.BaseModel `bun:"table:rates"`

	Id      string `bun:",pk,type:uuid"`
	UserId  string `bun:",notnull,unique:user_manga_idx,type:uuid"`
	MangaId string `bun:",notnull,unique:user_manga_idx,type:uuid"`
	Rate    uint8  `bun:",notnull"`

	CreatedAt time.Time `bun:",notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp"`

	User  *user_entity.User `bun:"rel:belongs-to,join:user_id=id"`
	Manga *Manga            `bun:"rel:belongs-to,join:manga_id=id,on_delete:CASCADE"`
}

func NewRate(userId, mangaId string, rate uint8) Rate {
	currentTime := time.Now()
	return Rate{
		Id:        uuid.NewString(),
		UserId:    userId,
		MangaId:   mangaId,
		Rate:      rate,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}
