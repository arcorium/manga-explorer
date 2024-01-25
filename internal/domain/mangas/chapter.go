package mangas

import (
	"github.com/biter777/countries"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"manga-explorer/internal/common"
	"manga-explorer/internal/domain/users"
	"time"
)

type Chapter struct {
	bun.BaseModel `bun:"table:chapters"`

	Id           string `bun:",pk,type:uuid,pk"`
	VolumeId     string `bun:",nullzero,notnull,type:uuid,unique:chapter_lang_idx"`
	TranslatorId string `bun:",nullzero,type:uuid,default:'afcd4ab0-3190-4d35-885a-1d20eb909bd9',unique:chapter_lang_idx"`

	Language    common.Language `bun:",notnull,unique:chapter_lang_idx,type:varchar(3)"`
	Title       string          `bun:",nullzero"`
	Number      uint64          `bun:",nullzero,notnull,unique:chapter_lang_idx"`
	PublishDate time.Time       `bun:",nullzero,type:date"`

	CreatedAt time.Time `bun:",notnull"`
	UpdatedAt time.Time `bun:",notnull"`

	Comments   []Comment   `bun:"rel:has-many,join:id=object_id,join:type=object_type,polymorphic"`
	Pages      []Page      `bun:"rel:has-many,join:id=chapter_id"`
	Translator *users.User `bun:"rel:belongs-to,join:translator_id=id,on_delete:SET DEFAULT"`
	Volume     *Volume     `bun:"rel:belongs-to,join:volume_id=id,on_delete:CASCADE"`
}

func NewChapter(volumeId, translatorId, title string, lang countries.CountryCode, number uint64, publishDate time.Time) Chapter {
	currentTime := time.Now()
	return Chapter{
		Id:           uuid.NewString(),
		VolumeId:     volumeId,
		TranslatorId: translatorId,
		Language:     common.Language(lang.Alpha3()),
		Title:        title,
		Number:       number,
		PublishDate:  publishDate,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
}
