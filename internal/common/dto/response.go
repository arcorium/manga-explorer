package dto

import (
	"manga-explorer/internal/common/status"
)

func NewErrorResponse(stat status.Object, details any) ErrorResponse {
	var detail any
	if details == nil && len(stat.DetailMessage()) > 0 {
		detail = stat.DetailMessage()
	} else {
		detail = details
	}
	return ErrorResponse{
		Code:    stat.Code,
		Message: stat.ErrorMessage(),
		Details: detail,
	}
}

func NewSuccessResponse(Code uint, data any, page *ResponsePage) SuccessWrapper {
	return SuccessWrapper{
		Internal: SuccessResponse{
			Code: Code,
			Data: data,
			Page: page,
		},
	}
}

type SuccessWrapper struct {
	Internal SuccessResponse `json:"success"`
}

type ErrorResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// MessageData could be used to substitute generic type on response as data
type MessageData struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Code uint          `json:"code"`
	Data any           `json:"data,omitempty"`
	Page *ResponsePage `json:"page,omitempty"`
}

type ErrorWrapper struct {
	Internal ErrorResponse `json:"error"`
}

type ResponsePage struct {
	Elements      uint64 `json:"elements"`
	CurrentPage   uint64 `json:"page,omitempty"`
	TotalElements uint64 `json:"total_elements"`
	TotalPage     uint64 `json:"total_page,omitempty"`
}
