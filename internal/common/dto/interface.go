package dto

import "manga-explorer/internal/common"

type AuthorizedDTO interface {
  SetUserId(claims *common.AccessTokenClaims)
}
