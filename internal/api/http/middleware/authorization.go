package middleware

import (
  "errors"
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt"
  "manga-explorer/internal/common"
  "manga-explorer/internal/common/status"
  "manga-explorer/internal/util"
  "manga-explorer/internal/util/httputil/resp"
  "strings"
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

func (a AuthMiddleware) getToken(ctx *gin.Context) (string, status.Object) {
  data := ctx.GetHeader(a.config.HeaderLookUp)
  if len(data) == 0 {
    return "", status.Error(status.AUTH_UNAUTHORIZED)
  }

  split := strings.Split(data, " ")
  if len(split) != 2 {
    return "", status.Error(status.TOKEN_LOOKUP_MALFORMED)
  }

  if split[0] != a.config.TokenType {
    return "", status.Error(status.TOKEN_MALTYPE)
  }

  return split[1], status.InternalSuccess()
}

func (a AuthMiddleware) getClaims(tokenStr string) (common.AccessTokenClaims, status.Object) {
  claims := common.AccessTokenClaims{}
  _, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(a.config.secretKey), nil
  })
  if err != nil {
    var verr *jwt.ValidationError
    ok := errors.As(err, &verr)
    if !ok {
      return common.AccessTokenClaims{}, status.Error(status.INTERNAL_SERVER_ERROR)
    }

    if verr.Errors&jwt.ValidationErrorExpired != 0 {
      return common.AccessTokenClaims{}, status.Error(status.ACCESS_TOKEN_EXPIRED)
    } else if verr.Errors&jwt.ValidationErrorMalformed != 0 {
      return common.AccessTokenClaims{}, status.Error(status.TOKEN_MALFORMED)
    }
  }

  // Response
  if err := claims.Valid(); err != nil {
    return common.AccessTokenClaims{}, status.Error(status.TOKEN_NOT_VALID)
  }

  return claims, status.InternalSuccess()
}

func (a AuthMiddleware) Handle(ctx *gin.Context) {
  // Get and validate token from header
  tokenStr, stat := a.getToken(ctx)
  if stat.IsError() {
    resp.Error(ctx, stat)
    ctx.Abort()
    return
  }

  // Parse
  claims, stat := a.getClaims(tokenStr)
  if stat.IsError() {
    resp.Error(ctx, stat)
    ctx.Abort()
    return
  }

  // Set claims on context
  ctx.Set(a.config.ClaimsKey, &claims)

  ctx.Next()
}

// Handle2 Works like Handle, but it will continue even when the authorization token is not found or malformed.
// Used when the controller can handle both and aware of it
func (a AuthMiddleware) Handle2(ctx *gin.Context) {
  // Get and validate token from header
  tokenStr, stat := a.getToken(ctx)
  if stat.IsError() {
    ctx.Next()
    return
  }

  // Parse
  claims, stat := a.getClaims(tokenStr)
  if stat.IsError() {
    ctx.Next()
    return
  }

  // Set claims on context
  ctx.Set(a.config.ClaimsKey, &claims)

  ctx.Next()
}
