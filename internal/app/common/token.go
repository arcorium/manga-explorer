package common

import (
	"github.com/golang-jwt/jwt"
)

type AccessTokenClaims struct {
	Id     string `json:"id"`
	UserId string `json:"uid"`
	Name   string `json:"name"`
	Role   string `json:"role"`

	jwt.StandardClaims
}
