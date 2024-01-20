package users

import (
	"github.com/uptrace/bun"
	"time"
)

func NewProfile(user *User, firstName, lastName string) Profile {
	return Profile{
		UserId:    user.Id,
		FirstName: firstName,
		LastName:  lastName,
	}
}

type Profile struct {
	bun.BaseModel `bun:"table:profiles"`

	Id        uint64 `bun:",pk,autoincrement"`
	UserId    string `bun:",notnull,unique,type:uuid"`
	FirstName string `bun:",notnull,nullzero"`
	LastName  string `bun:",nullzero"`
	PhotoURL  string `bun:",nullzero,type:text"`
	Bio       string `bun:",nullzero,type:text"`

	UpdatedAt time.Time `bun:",notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
}
