package httputil

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common"
	"manga-explorer/internal/app/dto"
)

type successWrapper[T any] struct {
	Internal T `json:"success"`
}

type errorWrapper[T any] struct {
	Internal T `json:"error"`
}

// SuccessResponse Used to set common.Response as response for success response
func SuccessResponse(ctx *gin.Context, status common.Status, data any) {
	res := successWrapper[dto.SuccessResponse]{Internal: dto.NewSuccessResponse(data, nil)}
	ctx.JSON(HttpCodeFromError(status), res)
}

func SuccessResponseMessage(ctx *gin.Context, status common.Status, message string) {
	res := successWrapper[dto.SuccessResponse]{Internal: dto.NewSuccessResponse(struct {
		Message string `json:"message"`
	}{message}, nil)}
	ctx.JSON(HttpCodeFromError(status), res)
}

func SuccessPagedResponse(ctx *gin.Context, status common.Status, page *dto.ResponsePage, data any) {
	res := successWrapper[dto.SuccessResponse]{Internal: dto.NewSuccessResponse(data, page)}
	ctx.JSON(HttpCodeFromError(status), res)
}

// ErrorResponse Used to set common.Response as response for bad response
func ErrorResponse[T any](ctx *gin.Context, status common.Status, details ...T) {
	res := errorWrapper[dto.ErrorResponse[T]]{Internal: dto.NewErrorResponse[T](status, details...)}
	ctx.JSON(HttpCodeFromError(status), res)
}

// Response Used to set response based on the common.Status passed as parameter. depending
// on the value of the common.Status it will call ErrorResponse or SuccessResponse
func Response(ctx *gin.Context, status common.Status, successData any) {
	if status.IsError() {
		ErrorResponse[string](ctx, status)
	} else {
		SuccessResponse(ctx, status, successData)
	}
}

func PagedResponse(ctx *gin.Context, status common.Status, successPage *dto.ResponsePage, successData any) {
	if status.IsError() {
		ErrorResponse[string](ctx, status)
	} else {
		SuccessPagedResponse(ctx, status, successPage, successData)
	}
}
