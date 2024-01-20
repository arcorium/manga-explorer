package common

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

var ErrUnknownRegion = errors.New("region unknown")

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func GetFieldsError(err validator.ValidationErrors) []FieldError {
	res := []FieldError{}
	for _, verr := range err {
		res = append(res, FieldError{
			Field: verr.Field(),
			Error: verr.Error(),
		})
	}
	return res
}
