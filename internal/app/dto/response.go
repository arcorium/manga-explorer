package dto

import (
	"manga-explorer/internal/app/common/status"
)

func NewErrorResponse(stat status.Object, details any) ErrorResponse {
	var detail any
	if details == nil && len(stat.DetailMessage()) > 0 {
		detail = stat.DetailMessage()
	}
	return ErrorResponse{
		Code:    stat.Code,
		Message: stat.ErrorMessage(),
		Details: detail,
	}
}

func NewSuccessResponse(data any, page *ResponsePage) SuccessResponse {
	return SuccessResponse{
		Data: data,
		Page: page,
	}
}

type ErrorResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
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
