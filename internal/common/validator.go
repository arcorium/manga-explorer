package common

import (
	"github.com/go-playground/validator/v10"
	"manga-explorer/internal/util"
)

func validateUUID4Slice(level validator.FieldLevel) bool {
	val, ok := level.Field().Interface().([]string)
	if !ok {
		return false
	}

	for _, v := range val {
		if !util.IsUUID(v) {
			return false
		}
	}
	return true
}

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

	err := validate.RegisterValidation("uuid4s", validateUUID4Slice)
	if err != nil {
		panic(err)
	}

	validate.RegisterAlias("manga_status", "oneof=completed ongoing drafted dropped hiatus")
}
