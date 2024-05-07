package common

import (
  "github.com/go-playground/validator/v10"
)

func RegisterValidationTags(validate *validator.Validate) {
  //err := validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
  //	if fl.Field().IsNil() {
  //		return false
  //	}
  //	return util.IsUUID(fl.Field().String())
  //})
  //
  //if err != nil {
  //	panic(err)
  //}

  //err := validate.RegisterValidation("language", validateLanguages)
  //if err != nil {
  //	panic(err)
  //}
  validate.RegisterAlias("language", "bcp47_language_tag")

  validate.RegisterAlias("manga_status", "oneof=completed ongoing drafted dropped hiatus")
}
