package dto

import (
	"manga-explorer/internal/app/common"
)

func NewErrorResponse[T any](stat common.Status, details ...T) ErrorResponse[T] {
	return ErrorResponse[T]{
		Code:    stat.Code,
		Message: stat.ErrorMessage(),
		Details: details,
	}
}

func NewSuccessResponse(data any, page *ResponsePage) SuccessResponse {
	return SuccessResponse{
		Data: data,
		Page: page,
	}
}

type ErrorResponse[T any] struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Details []T    `json:"details"`
}

type SuccessResponse struct {
	Data any           `json:"data,omitempty"`
	Page *ResponsePage `json:"page,omitempty"`
}

type ResponsePage struct {
	Elements      uint64 `json:"elements"`
	CurrentPage   uint64 `json:"page"`
	TotalElements uint64 `json:"total_elements"`
	TotalPage     uint64 `json:"total_page"`
}
