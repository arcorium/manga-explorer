package common

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var ErrUnknownRegion = errors.New("region unknown")
var ErrNilObject = errors.New("nil")

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ParameterError struct {
	Param string `json:"param"`
	Error string `json:"error"`
}

func GetFieldsError(err validator.ValidationErrors) []FieldError {
	res := []FieldError{}
	for _, verr := range err {
		errMessage := ""
		if verr.Tag() == "required" {
			errMessage = fmt.Sprintf("%s field is required", strings.ToLower(verr.Field()))
		} else {
			errMessage = fmt.Sprintf("%s field should satisfy %s constraint", strings.ToLower(verr.Field()), verr.Tag())
		}
		res = append(res, FieldError{
			Field: verr.Field(),
			Error: errMessage,
		})
	}
	return res
}

func NewNotPresentParameter(paramName string) ParameterError {
	return ParameterError{
		Param: paramName,
		Error: "parameter is empty",
	}
}

func NewParameterError(paramName string, errMessage string) ParameterError {
	return ParameterError{
		Param: paramName,
		Error: fmt.Sprintf("%s parameter %s", paramName, errMessage),
	}
}
