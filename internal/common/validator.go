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

	validate.RegisterAlias("manga_status", "oneof=completed ongoing drafted dropped hiatus")
}
