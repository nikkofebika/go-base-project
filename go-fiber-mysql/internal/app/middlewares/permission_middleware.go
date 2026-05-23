package middlewares

import (
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/enums"
	"go-fiber-mysql/internal/app/exception"
	"go-fiber-mysql/internal/app/services"

	"github.com/gofiber/fiber/v3"
)

const UserPermissionsContextKey = "user_permissions"

func PermissionMiddleware(service services.UserService, requiredPermissions ...enums.Permission) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		// 1. Get user from context (set by AuthMiddleware)
		userVal := ctx.Locals("user")
		if userVal == nil {
			return exception.NewUnauthorizedException()
		}

		user, ok := userVal.(*entities.User)
		if !ok {
			return exception.NewUnauthorizedException()
		}

		// Bypass for admin
		if user.Type == enums.UserTypeAdmin {
			return ctx.Next()
		}

		userID := user.ID

		// 2. Check for cached permissions in Locals
		var userPermissions []string
		cachedPermissions := ctx.Locals(UserPermissionsContextKey)

		if cachedPermissions != nil {
			userPermissions, _ = cachedPermissions.([]string)
		} else {
			// 3. Fetch from DB if not cached
			var err error
			userPermissions, err = service.GetPermissionSlugsByUserID(ctx.Context(), userID)
			if err != nil {
				return exception.NewInternalServerException(err.Error())
			}

			// 4. Cache in Locals for this request
			ctx.Locals(UserPermissionsContextKey, userPermissions)
		}

		// 5. Check if user has ANY of the required permissions
		hasPermission := false
		for _, required := range requiredPermissions {
			for _, userPerm := range userPermissions {
				if userPerm == required.String() {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			return exception.NewForbiddenException()
		}

		return ctx.Next()
	}
}
