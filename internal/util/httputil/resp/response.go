package resp

import (
	"github.com/gin-gonic/gin"
	"manga-explorer/internal/common/dto"
	"manga-explorer/internal/common/status"
)

// Success Used to set common.Response as response for success response
func Success(ctx *gin.Context, status status.Object, data any, page *dto.ResponsePage) {
	res := dto.NewSuccessResponse(status.Code, data, page)
	ctx.JSON(HttpCodeFromError(status), res)
}

func SuccessMessage(ctx *gin.Context, status status.Object, message string) {
	res := dto.NewSuccessResponse(status.Code, dto.MessageData{Message: message}, nil)
	ctx.JSON(HttpCodeFromError(status), res)
}

// ErrorDetailed Used to set common.Response as response for bad response
func ErrorDetailed(ctx *gin.Context, status status.Object, details any) {
	res := dto.ErrorWrapper{Internal: dto.NewErrorResponse(status, details)}
	ctx.JSON(HttpCodeFromError(status), res)
}

func Error(ctx *gin.Context, status status.Object) {
	res := dto.ErrorWrapper{Internal: dto.NewErrorResponse(status, nil)}
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
