package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func existsValidator(db *gorm.DB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		parts := strings.Split(fl.Param(), ".")
		if len(parts) != 2 {
			return false
		}

		table := parts[0]
		column := parts[1]

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

		var value any

		switch field.Kind() {
		case reflect.Uint, reflect.Uint64, reflect.Uint32:
			value = field.Uint()
		case reflect.Int, reflect.Int64, reflect.Int32:
			value = field.Int()
		case reflect.String:
			value = field.String()
		default:
			return false
		}

		var count int64
		if err := db.Table(table).
			Where(column+" = ?", value).
			Limit(1).
			Count(&count).Error; err != nil {
			return false
		}

		return count > 0
	}
}

// func existsValidator(db *gorm.DB) validator.Func {
// 	return func(field validator.FieldLevel) bool {
// 		// Tag exists:table,column
// 		parts := strings.Split(field.Param(), ".")
// 		if len(parts) != 2 {
// 			return false
// 		}

// 		table := parts[0]
// 		column := parts[1]
// 		fieldValue := field.Field().Interface()

// 		var count int64
// 		if err := db.Table(table).Where(column+" = ?", fieldValue).Limit(1).Count(&count).Error; err != nil {
// 			return false
// 		}

// 		return count > 0
// 	}
// }
