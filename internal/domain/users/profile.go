package users

import (
  "github.com/uptrace/bun"
  "manga-explorer/internal/infrastructure/file"
  "manga-explorer/internal/util/opt"
  "time"
)

func NewProfile(userId, firstName, lastName, bio string, photoUrl opt.Optional[file.Name]) Profile {
  return Profile{
    UserId:    userId,
    FirstName: firstName,
    LastName:  lastName,
    PhotoURL:  photoUrl.ValueOr(""),
    Bio:       bio,
    UpdatedAt: time.Now(),
  }
}

type Profile struct {
  bun.BaseModel `bun:"table:profiles"`

  Id        uint64    `bun:",pk,autoincrement"`
  UserId    string    `bun:",nullzero,notnull,unique,type:uuid"`
  FirstName string    `bun:",nullzero,notnull"`
  LastName  string    `bun:","`
  PhotoURL  file.Name `bun:",type:text"`
  Bio       string    `bun:",type:text"`

  UpdatedAt time.Time `bun:",nullzero,notnull"`

  User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
}
