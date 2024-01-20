package httputil

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/app/dto"
	"manga-explorer/internal/util"
)

func BindUriJson[T any](ctx *gin.Context) (T, common.Status) {
	var t T
	if err := ctx.BindUri(&t); err != nil {
		return t, common.StatusError(status.BAD_PARAMETER_ERROR)
	}

	if err := ctx.BindJSON(&t); err != nil {
		return t, common.StatusError(status.BAD_BODY_REQUEST_ERROR)
	}
	return t, common.StatusSuccess()
}

func BindQueryJson[T any](ctx *gin.Context) (T, common.Status) {
	var t T
	if err := ctx.BindQuery(&t); err != nil {
		return t, common.StatusError(status.BAD_QUERY_ERROR)
	}

	if err := ctx.BindJSON(&t); err != nil {
		return t, common.StatusError(status.BAD_BODY_REQUEST_ERROR)
	}
	return t, common.StatusSuccess()
}

func BindUriMultipartForm[T any](ctx *gin.Context) (T, common.Status) {
	var t T
	if err := ctx.BindUri(&t); err != nil {
		return t, common.StatusError(status.BAD_QUERY_ERROR)
	}

	if err := ctx.Bind(&t); err != nil {
		return t, common.StatusError(status.BAD_BODY_REQUEST_ERROR)
	}
	return t, common.StatusSuccess()
}

func BindAuthorizedJSON(ctx *gin.Context, data dto.AuthorizedDTO) common.Status {
	claims, err := util.GetContextValue[*common.AccessTokenClaims](ctx, common.ClaimsKey)
	if err != nil {
		return common.StatusError(status.AUTH_UNAUTHORIZED)
	}

	if err = ctx.BindJSON(&data); err != nil {
		return common.StatusError(status.BAD_BODY_REQUEST_ERROR)
	}
	data.SetUserId(claims)
	return common.StatusSuccess()
}
