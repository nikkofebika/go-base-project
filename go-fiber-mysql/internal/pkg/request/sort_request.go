package request

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type SortFunc func(db *gorm.DB, direction string) *gorm.DB

type AllowedSorts map[string]SortFunc

type AppliedSort struct {
	Field     string
	Apply     SortFunc
	Direction string
}

func ParseSort(
	ctx fiber.Ctx,
	allowed AllowedSorts,
) []AppliedSort {

	var sorts []AppliedSort

	sortParam := ctx.Query("sort")
	if sortParam == "" {
		return sorts
	}

	// Split multiple sorts if comma-separated
	sortFields := strings.Split(sortParam, ",")

	for _, sortField := range sortFields {
		sortField = strings.TrimSpace(sortField)

		// Determine direction (default ascending)
		direction := "ASC"
		field := sortField

		if strings.HasPrefix(sortField, "-") {
			direction = "DESC"
			field = sortField[1:]
		}

		// Check if field is allowed
		if fn, ok := allowed[field]; ok {
			sorts = append(sorts, AppliedSort{
				Field:     field,
				Apply:     fn,
				Direction: direction,
			})
		}
	}

	return sorts
}

func ApplySorts(db *gorm.DB, sorts []AppliedSort) *gorm.DB {
	for _, sort := range sorts {
		db = sort.Apply(db, sort.Direction)
	}
	return db
}

// Helper function for simple column sorting
func SortByColumn(column string) SortFunc {
	return func(db *gorm.DB, direction string) *gorm.DB {
		return db.Order(column + " " + direction)
	}
}

// Helper function for sorting by related table column with join
func SortByRelation(joinClause string, column string) SortFunc {
	return func(db *gorm.DB, direction string) *gorm.DB {
		return db.Joins(joinClause).Order(column + " " + direction)
	}
}
