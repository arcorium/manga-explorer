package users

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func NewCredential(user *User, deviceName, accessTokenId, token string) Credential {
	return Credential{
		Id:            uuid.NewString(),
		UserId:        user.Id,
		AccessTokenId: accessTokenId,
		Device:        Device{Name: deviceName},
		Token:         token,
	}
}

func NewCredential2(userId, deviceName, accessTokenId, token string) Credential {
	return Credential{
		BaseModel:     bun.BaseModel{},
		Id:            uuid.NewString(),
		UserId:        userId,
		AccessTokenId: accessTokenId,
		Device: Device{
			Name: deviceName,
		},
		Token:     token,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

type Credential struct {
	bun.BaseModel `bun:"table:credentials"`

	Id            string `bun:",pk,type:uuid"`
	UserId        string `bun:",notnull,type:uuid"`
	AccessTokenId string `bun:",notnull,type:uuid"` // used to prevent creating many access tokens from one expired access token
	Device        Device `bun:"embed:device_"`
	Token         string `bun:",notnull"`

	UpdatedAt time.Time `bun:",default:current_timestamp,notnull"`
	CreatedAt time.Time `bun:",default:current_timestamp,notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
}
