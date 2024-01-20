package dto

import "manga-explorer/internal/app/common"

type AuthorizedDTO interface {
	SetUserId(claims *common.AccessTokenClaims)
}
