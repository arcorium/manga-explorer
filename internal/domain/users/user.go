package users

import (
	"manga-explorer/internal/app/common/constant"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"manga-explorer/internal/util"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9|._]+@[a-zA-Z0-9|-]+(\\.[a-zA-Z0-9]{2,})+$")

func NewUser(name, email, password string, role Role) (User, error) {
	currentTime := time.Now()
	usr := User{
		Id:        uuid.NewString(),
		Username:  name,
		Email:     email,
		Verified:  false,
		Role:      role,
		UpdatedAt: currentTime,
		CreatedAt: currentTime,
	}

	password, err := util.Hash(password)
	if err != nil {
		return BadUser, ErrHashPassword
	}
	usr.Password = password

	if !usr.ValidateEmail() {
		return BadUser, ErrEmailValidation
	}
	return usr, nil
}

var BadUser User

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id       string `bun:",type:uuid,pk"`
	Username string `bun:",notnull,nullzero,unique"`
	Email    string `bun:",notnull,nullzero,unique,type:text"`
	Password string `bun:",notnull,nullzero"`
	Verified bool   `bun:"is_verified,default:false"`
	Role     Role   `bun:",notnull"`

	BannedUntil time.Time `bun:",nullzero,default:null"`
	DeletedAt   time.Time `bun:",nullzero,default:null"`
	UpdatedAt   time.Time `bun:",notnull"`
	CreatedAt   time.Time `bun:",notnull"`
}

func (u *User) ValidatePassword(rawPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rawPassword)) == nil
}

func (u *User) GenerateAccessTokenClaims(duration time.Duration) jwt.MapClaims {
	defaultClaims := DefaultClaims(duration, constant.IssuerName)
	defaultClaims["id"] = uuid.NewString()
	defaultClaims["uid"] = u.Id
	defaultClaims["name"] = u.Username
	defaultClaims["role"] = u.Role.String() // TODO: Maybe better to have the uint8 instead
	return defaultClaims
}

func (u *User) ValidateEmail() bool {
	return emailRegex.MatchString(u.Email)
}

func DefaultClaims(duration time.Duration, issuer string) jwt.MapClaims {
	currentTime := time.Now()
	return jwt.MapClaims{
		"iat": currentTime.Unix(),
		"nbf": currentTime.Unix(),
		"iss": issuer,
		"exp": currentTime.Add(duration).Unix(),
	}
}
