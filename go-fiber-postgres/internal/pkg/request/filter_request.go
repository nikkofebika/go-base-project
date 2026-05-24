package request

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type FilterFunc func(db *gorm.DB, value string) *gorm.DB

type AppliedFilter struct {
	Field string
	Apply FilterFunc
	Value string
}

type AllowedFilters map[string]FilterFunc

func ParseFilters(ctx fiber.Ctx, allowedFilters AllowedFilters) []AppliedFilter {
	var filters []AppliedFilter

	queries := ctx.Queries()
	for key, value := range queries {
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			field := key[7 : len(key)-1]

			if fn, ok := allowedFilters[field]; ok {
				filters = append(filters, AppliedFilter{
					Field: field,
					Apply: fn,
					Value: value,
				})
			}
		}
	}

	return filters
}

func FilterExact(column string) FilterFunc {
	return func(db *gorm.DB, value string) *gorm.DB {
		values := strings.Split(value, ",")

		if len(values) == 1 {
			return db.Where(column+" = ?", value)
		}

		// Remove empty values and trim whitespace
		filteredValues := make([]string, 0, len(values))
		for _, v := range values {
			v = strings.TrimSpace(v)
			if v != "" {
				filteredValues = append(filteredValues, v)
			}
		}

		if len(filteredValues) == 0 {
			return db
		}

		if len(filteredValues) == 1 {
			return db.Where(column+" = ?", filteredValues[0])
		}

		// Use IN clause for multiple values
		placeholders := make([]string, len(filteredValues))
		args := make([]any, len(filteredValues))
		for i, v := range filteredValues {
			placeholders[i] = "?"
			args[i] = v
		}

		return db.Where(column+" IN ("+strings.Join(placeholders, ", ")+")", args...)
	}
}

func FilterLike(columns ...string) FilterFunc {
	return func(db *gorm.DB, value string) *gorm.DB {
		// if len(columns) > 0 {
		conditions := make([]string, 0, len(columns))
		args := make([]any, 0, len(columns))

		for _, col := range columns {
			conditions = append(conditions, col+" LIKE ?")
			args = append(args, "%"+value+"%")
		}

		return db.Where(strings.Join(conditions, " OR "), args...)
		// }
	}
}
