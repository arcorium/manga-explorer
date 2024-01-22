package resp

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/app/common/status"
	"manga-explorer/internal/app/dto"
)

//type SuccessWrapper[T any] struct {
//	Internal T `json:"success"`
//}

type ErrorWrapper[T any] struct {
	Internal T `json:"error"`
}

// Success Used to set common.Response as response for success response
func Success(ctx *gin.Context, status status.Object, data any, page *dto.ResponsePage) {
	//res := SuccessWrapper[dto.Success]{Internal: dto.NewSuccessResponse(data, page)}
	ctx.JSON(HttpCodeFromError(status), dto.NewSuccessResponse(data, page))
}

func SuccessMessage(ctx *gin.Context, status status.Object, message string) {
	//res := SuccessWrapper[dto.Success]{Internal: dto.NewSuccessResponse(struct {
	//	Message string `json:"message"`
	//}{message}, nil)}
	ctx.JSON(HttpCodeFromError(status), dto.NewSuccessResponse(struct {
		Message string `json:"message"`
	}{message}, nil))
}

// ErrorDetailed Used to set common.Response as response for bad response
func ErrorDetailed[T any](ctx *gin.Context, status status.Object, details T) {
	res := ErrorWrapper[dto.ErrorResponse]{Internal: dto.NewErrorResponse(status, details)}
	ctx.JSON(HttpCodeFromError(status), res)
}

func Error(ctx *gin.Context, status status.Object) {
	res := ErrorWrapper[dto.ErrorResponse]{Internal: dto.NewErrorResponse(status, nil)}
	ctx.JSON(HttpCodeFromError(status), res)

}

// Conditional Used to set response based on the common.Object passed as parameter. depending
// on the value of the common.Object it will call Error or Success
func Conditional(ctx *gin.Context, status status.Object, successData any, page *dto.ResponsePage) {
	if status.IsError() {
		Error(ctx, status)
	} else {
		Success(ctx, status, successData, page)
	}
}
