package validator

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func LengthValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		// 1. Split the param "1.100" into ["1", "100"]
		parts := strings.Split(fl.Param(), ".")
		if len(parts) != 2 {
			return false
		}

		min, errMin := strconv.Atoi(parts[0])
		max, errMax := strconv.Atoi(parts[1])
		if errMin != nil || errMax != nil {
			return false
		}

		field := fl.Field()

		// ===== HANDLE Nullable[T] =====
		if field.Kind() == reflect.Struct {
			setField := field.FieldByName("Set")
			valueField := field.FieldByName("Value")

			// Bukan Nullable struct
			if setField.IsValid() && valueField.IsValid() {
				// Tidak dikirim → skip
				if !setField.Bool() {
					return true
				}

				// Dikirim null → skip exists
				if valueField.IsNil() {
					return true
				}

				// Ambil value aslinya
				field = valueField.Elem()
			}
		}

		// ===== HANDLE pointer biasa =====
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return true
			}
			field = field.Elem()
		}

		var value int

		switch field.Kind() {
		case reflect.Uint, reflect.Uint64, reflect.Uint32:
			value = int(field.Uint())
		case reflect.Int, reflect.Int64, reflect.Int32:
			value = int(field.Int())
		case reflect.String:
			value = len(field.String())
		default:
			return false
		}

		return value >= min && value <= max
	}
}
