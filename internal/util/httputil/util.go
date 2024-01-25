package httputil

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"manga-explorer/internal/common"
	"manga-explorer/internal/common/status"
)

func BindJson[T any](ctx *gin.Context, data *T) (status.Object, []common.FieldError) {
	if err := ctx.BindJSON(data); err != nil {
		var verr validator.ValidationErrors
		errors.As(err, &verr)
		return status.Error(status.BAD_REQUEST_ERROR), common.GetFieldsError(verr)
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

func BindMultipartForm[T any](ctx *gin.Context, data *T) (status.Object, []common.FieldError) {
	if err := ctx.Bind(data); err != nil {
		var verr validator.ValidationErrors
		errors.As(err, &verr)
		return status.Error(status.BAD_REQUEST_ERROR), common.GetFieldsError(verr)
	}
	return status.Success(), nil
}
