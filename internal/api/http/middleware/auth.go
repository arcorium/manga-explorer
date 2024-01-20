package middleware

import (
	"errors"
	"log"
	"manga-explorer/internal/app/common/status"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/util"
	"manga-explorer/internal/util/httputil"
)

type AuthMiddlewareConfig struct {
	SigningMethod jwt.SigningMethod // Default: HS256
	secretKey     string

	TokenType    string // Default: Bearer
	ClaimsKey    string // Default: claims
	HeaderLookUp string // Default: Authorization
}

func NewAuthMiddleware(secretKey string, config *AuthMiddlewareConfig) AuthMiddleware {
	if config == nil {
		config = &AuthMiddlewareConfig{}
	}

	// Set default config
	config.secretKey = secretKey
	util.SetDefaultString(&config.HeaderLookUp, "Authorization")
	util.SetDefaultString(&config.TokenType, "Bearer")
	util.SetDefaultString(&config.ClaimsKey, "claims")
	if config.SigningMethod == nil {
		config.SigningMethod = jwt.SigningMethodHS256
	}

	return AuthMiddleware{config: config}
}

type AuthMiddleware struct {
	config *AuthMiddlewareConfig
}

func (a AuthMiddleware) Handle(ctx *gin.Context) {
	data := ctx.GetHeader(a.config.HeaderLookUp)
	if len(data) == 0 {
		httputil.ErrorResponse(ctx, common.StatusError(status.AUTH_UNAUTHORIZED))
		ctx.Abort()
		return
	}

	split := strings.Split(data, " ")
	if len(split) != 2 {
		httputil.ErrorResponse(ctx, common.StatusError(status.TOKEN_LOOKUP_MALFORMED))
		ctx.Abort()
		return
	}

	if split[0] != a.config.TokenType {
		httputil.ErrorResponse(ctx, common.StatusError(status.TOKEN_MALTYPE))
		ctx.Abort()
		return
	}

	// Parse
	claims := common.AccessTokenClaims{}
	_, err := jwt.ParseWithClaims(split[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.secretKey), nil
	})
	if err != nil {
		var verr *jwt.ValidationError
		ok := errors.As(err, &verr)
		if !ok {
			httputil.ErrorResponse(ctx, common.StatusError(status.INTERNAL_SERVER_ERROR))
			ctx.Abort()
			return
		}

		if verr.Errors&jwt.ValidationErrorExpired != 0 {
			httputil.ErrorResponse(ctx, common.StatusError(status.ACCESS_TOKEN_EXPIRED))
			ctx.Abort()
			return
		} else if verr.Errors&jwt.ValidationErrorMalformed != 0 {
			httputil.ErrorResponse(ctx, common.StatusError(status.TOKEN_MALFORMED))
			ctx.Abort()
			return
		}
	}

	// Validate
	if err := claims.Valid(); err != nil {
		log.Println(err)
		httputil.ErrorResponse(ctx, common.StatusError(status.TOKEN_NOT_VALID))
		ctx.Abort()
		return
	}

	log.Println(time.Unix(claims.ExpiresAt, 0))

	// Set claims on context
	ctx.Set(a.config.ClaimsKey, &claims)

	ctx.Next()
}
