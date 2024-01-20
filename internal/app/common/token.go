package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"manga-explorer/internal/app/common/constant"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/util"
)

type AccessTokenClaims struct {
	Id     string `json:"id"`
	UserId string `json:"uid"`
	Name   string `json:"name"`
	Role   string `json:"role"`

	jwt.StandardClaims
}

func GetClaims(ctx *gin.Context) (*AccessTokenClaims, status.Object) {
	value, err := util.GetContextValue[*AccessTokenClaims](ctx, constant.ClaimsKey)
	if errors.Is(err, util.NoContextValueErr) {
		return nil, status.Error(status.AUTH_UNAUTHORIZED)
	} else if errors.Is(err, util.MaltypeContextValueErr) {
		return nil, status.Error(status.TOKEN_MALFORMED)
	}
	return value, status.Success()
}
