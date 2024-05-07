package mangas

import (
  "github.com/uptrace/bun"
  "manga-explorer/internal/domain/users"
  "manga-explorer/internal/util/opt"
  "time"
)

type ChapterHistory struct {
  bun.BaseModel `bun:"table:chapter_histories"`

  UserId    string    `bun:",type:uuid,pk"`
  ChapterId string    `bun:",type:uuid,pk"`
  LastView  time.Time `bun:",nullzero,notnull"`

  User    *users.User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
  Chapter *Chapter    `bun:"rel:belongs-to,join:chapter_id=id,on_delete:CASCADE"`
}

func NewChapterHistory(userId, chapterId string, lastView opt.Optional[time.Time]) ChapterHistory {
  return ChapterHistory{
    UserId:    userId,
    ChapterId: chapterId,
    LastView:  lastView.ValueOr(time.Now()),
  }
}

// MangaHistory Used only when scanning or select from persistent storage
type MangaHistory struct {
  LastView time.Time `bun:",scanonly"`
  Manga    *Manga    `bun:",scanonly"`
}
