package mangas

import (
	"github.com/biter777/countries"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/infrastructure/file"
	"time"
)

type Manga struct {
	bun.BaseModel       `bun:"table:mangas"`
	Id                  string         `bun:",pk,type:uuid"`
	Status              Status         `bun:",nullzero,notnull"`
	Origin              common.Country `bun:",nullzero,notnull,type:varchar(2)"`
	OriginalTitle       string         `bun:",notnull,nullzero,unique,type:text"`
	OriginalDescription string         `bun:",nullzero,type:text"`
	PublicationYear     uint16         `bun:",notnull,nullzero"`
	CoverURL            file.Name      `bun:",nullzero"`
	UpdatedAt           time.Time      `bun:",nullzero,notnull"`
	CreatedAt           time.Time      `bun:",nullzero,notnull"`

	Comments     []Comment     `bun:"rel:has-many,join:id=object_id,join:type=object_type,polymorphic"`
	Ratings      []Rate        `bun:"rel:has-many,join:id=manga_id"`
	Translations []Translation `bun:"rel:has-many,join:id=manga_id"`
	Volumes      []Volume      `bun:"rel:has-many,join:id=manga_id"`
	Genres       []Genre       `bun:"m2m:manga_genres,join:Manga=Genre"`
}

func NewManga(title, desc, coverUrl string, year uint16, status Status, region countries.CountryCode) Manga {
	currentTime := time.Now()
	return Manga{
		Id:                  uuid.NewString(),
		Status:              status,
		Origin:              common.Country(region.Alpha2()),
		OriginalTitle:       title,
		OriginalDescription: desc,
		PublicationYear:     year,
		CoverURL:            file.Name(coverUrl),
		UpdatedAt:           currentTime,
		CreatedAt:           currentTime,
	}
}
