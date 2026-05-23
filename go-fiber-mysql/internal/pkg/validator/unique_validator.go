package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func uniqueValidator(db *gorm.DB) validator.Func {
	return func(field validator.FieldLevel) bool {
		// Tag unique:table,column
		parts := strings.Split(field.Param(), ".")
		if len(parts) != 2 {
			return false
		}

		table := parts[0]
		column := parts[1]
		fieldValue := field.Field().Interface()

		var count int64
		if err := db.Table(table).Where(column+" = ?", fieldValue).Limit(1).Count(&count).Error; err != nil {
			return false
		}

		return count == 0
	}
}
