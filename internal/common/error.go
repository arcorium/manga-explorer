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

    switch verr.Tag() {
    case "required":
      errMessage = fmt.Sprintf("%s field is required", strings.ToLower(verr.Field()))
    case "gt":
      errMessage = fmt.Sprintf("%s field length should be more than %s", strings.ToLower(verr.Field()), verr.Param())
    case "gte":
      errMessage = fmt.Sprintf("%s field length should be more than equal %s", strings.ToLower(verr.Field()), verr.Param())
    case "lt":
      errMessage = fmt.Sprintf("%s field length should be less than %s", strings.ToLower(verr.Field()), verr.Param())
    case "lte":
      errMessage = fmt.Sprintf("%s field length should be less than equal %s", strings.ToLower(verr.Field()), verr.Param())
    case "eq":
      errMessage = fmt.Sprintf("%s field length should be equal to %s", strings.ToLower(verr.Field()), verr.Param())
    case "ne":
      errMessage = fmt.Sprintf("%s field length should be not equal to %s", strings.ToLower(verr.Field()), verr.Param())
    default:
      errMessage = fmt.Sprintf("%s field should satisfy %s constraint", strings.ToLower(verr.Field()), strings.ReplaceAll(verr.Tag(), "|", " or "))
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
