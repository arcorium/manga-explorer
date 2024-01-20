package users

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func NewVerification(userId string, usage Usage) Verification {
	// Generate token with UUID as Underlying data
	id := uuid.NewString()
	hasher := sha256.New()
	hasher.Write([]byte(id))
	token := hex.EncodeToString(hasher.Sum(nil))

	return Verification{
		Token:          token,
		UserId:         userId,
		Usage:          usage,
		ExpirationTime: time.Now().Add(time.Minute * 15),
	}
}

type Verification struct {
	bun.BaseModel `bun:"table:verifications"`

	Token          string    `bun:",pk"`
	UserId         string    `bun:",notnull,type:uuid,unique:verif_usage_idx"`
	Usage          Usage     `bun:",unique:verif_usage_idx"`
	ExpirationTime time.Time `bun:",notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
}

type Usage uint8

func (u Usage) Underlying() uint8 {
	return uint8(u)
}

func (u Usage) Validate() error {
	if u.Underlying() > 1 {
		return ErrUnknownVerificationUsage
	}
	return nil
}

func (u Usage) String() string {
	switch u.Underlying() {
	case 0:
		return "reset password"
	case 1:
		return "email confirmation"
	}

	return "unknown"
}

const UsageResetPassword = Usage(0)
const UsageVerifyEmail = Usage(1)
