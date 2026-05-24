package validator

import (
	"go-fiber-postgres/internal/app/enums"

	"github.com/go-playground/validator/v10"
)

func enumValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		enum, ok := fl.Field().Interface().(enums.Enum)
		if !ok {
			return false
		}

		return enum.IsValid()
	}
}
