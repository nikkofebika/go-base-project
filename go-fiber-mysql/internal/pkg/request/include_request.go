package request

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type IncludeFunc func(db *gorm.DB) *gorm.DB

type AllowedIncludes map[string]IncludeFunc

type AppliedInclude struct {
	Field string
	Apply IncludeFunc
}

func ParseInclude(ctx fiber.Ctx, allowed AllowedIncludes) []AppliedInclude {
	var includes []AppliedInclude

	include := ctx.Query("include")
	if include == "" {
		return includes
	}

	for _, inc := range strings.Split(include, ",") {
		if fn, ok := allowed[inc]; ok {
			includes = append(includes, AppliedInclude{
				Field: inc,
				Apply: fn,
			})
		}
	}

	return includes
}

func IncludePreload(relation string, columns ...string) IncludeFunc {
	return func(db *gorm.DB) *gorm.DB {
		if len(columns) > 0 {
			return db.Preload(relation, func(db *gorm.DB) *gorm.DB {
				return db.Select(columns)
			})
		}

		return db.Preload(relation)
	}
}
