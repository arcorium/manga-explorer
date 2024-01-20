package httputil

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/common/constant"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/app/dto"
	"manga-explorer/internal/util"
)

func BindUriJson[T any](ctx *gin.Context, data *T) (status.Object, []common.FieldError) {
	if err := ctx.BindUri(data); err != nil {
		var verr validator.ValidationErrors
		errors.As(err, &verr)
		return status.Error(status.BAD_PARAMETER_ERROR), common.GetFieldsError(verr)
	}

	return BindJson(ctx, data)
}

func BindQueryJson[T any](ctx *gin.Context, data *T) (status.Object, []common.FieldError) {
	var t T
	if err := ctx.BindQuery(&t); err != nil {
		var verr validator.ValidationErrors
		errors.As(err, &verr)
		return status.Error(status.BAD_QUERY_ERROR), common.GetFieldsError(verr)
	}

	return BindJson(ctx, data)
}

func BindJson[T any](ctx *gin.Context, data *T) (status.Object, []common.FieldError) {
	if err := ctx.BindJSON(data); err != nil {
		var verr validator.ValidationErrors
		errors.As(err, &verr)
		return status.Error(status.BAD_BODY_REQUEST_ERROR), common.GetFieldsError(verr)
	}
	return status.Success(), nil
}

func BindUri[T any](ctx *gin.Context, data *T) (status.Object, []common.FieldError) {
	if err := ctx.BindUri(data); err != nil {
		var verr validator.ValidationErrors
		errors.As(err, &verr)
		return status.Error(status.BAD_PARAMETER_ERROR), common.GetFieldsError(verr)
	}
	return status.Success(), nil
}

func BindQuery[T any](ctx *gin.Context, data *T) (status.Object, []common.FieldError) {
	if err := ctx.BindQuery(data); err != nil {
		var verr validator.ValidationErrors
		errors.As(err, &verr)
		return status.Error(status.BAD_QUERY_ERROR), common.GetFieldsError(verr)
	}
	return status.Success(), nil
}

func BindUriMultipartForm[T any](ctx *gin.Context, data *T) (status.Object, []common.FieldError) {
	if err := ctx.BindUri(data); err != nil {
		var verr validator.ValidationErrors
		errors.As(err, &verr)
		return status.Error(status.BAD_PARAMETER_ERROR), common.GetFieldsError(verr)
	}

	if err := ctx.Bind(data); err != nil {
		var verr validator.ValidationErrors
		errors.As(err, &verr)
		return status.Error(status.BAD_BODY_REQUEST_ERROR), common.GetFieldsError(verr)
	}
	return status.Success(), nil
}

func BindAuthorizedJSON(ctx *gin.Context, data dto.AuthorizedDTO) (status.Object, []common.FieldError) {
	claims, err := util.GetContextValue[*common.AccessTokenClaims](ctx, constant.ClaimsKey)
	if err != nil {
		return status.Error(status.AUTH_UNAUTHORIZED), nil
	}

	if err = ctx.BindJSON(&data); err != nil {
		var verr validator.ValidationErrors
		errors.As(err, &verr)
		return status.Error(status.BAD_BODY_REQUEST_ERROR), common.GetFieldsError(verr)
	}
	data.SetUserId(claims)
	return status.Success(), nil
}
