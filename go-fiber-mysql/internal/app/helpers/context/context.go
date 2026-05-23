package context

import (
	"context"
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/exception"

	"github.com/gofiber/fiber/v3"
)

const (
	UserIDContextKey = "user_id"
	UserContextKey   = "user"
)

// ExtractUserIDFromLocals safely extracts user_id from fiber context locals
// Returns the user_id and an error if extraction fails
func ExtractUserIDFromLocals(ctx fiber.Ctx) (uint, error) {
	userIDVal := ctx.Locals(UserIDContextKey)
	if userIDVal == nil {
		return 0, exception.NewUnauthorizedException()
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		return 0, exception.NewUnauthorizedException()
	}

	return userID, nil
}

// ContextWithUserID creates a new context with user_id attached
// This replaces the repetitive pattern of extracting, converting, and setting user_id
func ContextWithUserID(ctx fiber.Ctx) (context.Context, uint, error) {
	userID, err := ExtractUserIDFromLocals(ctx)
	if err != nil {
		return nil, 0, err
	}

	// c := ctx.Context()
	// c.SetUserValue(UserIDContextKey, userID)
	context := context.WithValue(ctx.Context(), UserIDContextKey, userID)

	return context, userID, nil
}

// ContextWithUser creates a new context with user attached
// This replaces the repetitive pattern of extracting, converting, and setting user
func ContextWithUser(ctx fiber.Ctx) (context.Context, *entities.User, error) {
	user, err := ExtractUserFromLocals(ctx)
	if err != nil {
		return nil, nil, err
	}

	// c := ctx.Context()
	// c.SetUserValue(UserContextKey, user)
	context := context.WithValue(ctx.Context(), UserContextKey, user)

	return context, user, nil
}

func ExtractUserFromLocals(ctx fiber.Ctx) (*entities.User, error) {
	userVal := ctx.Locals(UserContextKey)
	if userVal == nil {
		return nil, exception.NewUnauthorizedException()
	}

	user, ok := userVal.(*entities.User)
	if !ok {
		return nil, exception.NewUnauthorizedException()
	}

	return user, nil
}

// ValidateParamID validates an ID parameter is a positive integer
// Returns the ID as uint and an error if validation fails
func ValidateParamID(ctx fiber.Ctx, paramName string) (uint, error) {
	// id, err := ctx.ParamsInt(paramName)
	id := fiber.Params[int](ctx, paramName)
	// if err != nil {
	// 	return 0, exception.NewBadRequestException("Invalid " + paramName + " format")
	// }

	if id <= 0 {
		return 0, exception.NewBadRequestException(paramName + " must be greater than 0")
	}

	return uint(id), nil
}

// ExtractUserIDFromContext safely extracts user_id from a context
// Used by service layer to get user_id that was set by handler
func ExtractUserIDFromContext(ctx context.Context) (uint, error) {
	userIDVal := ctx.Value(UserIDContextKey)
	if userIDVal == nil {
		return 0, exception.NewUnauthorizedException()
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		return 0, exception.NewUnauthorizedException()
	}

	return userID, nil
}

func ExtractUserFromContext(ctx context.Context) (*entities.User, error) {
	userVal := ctx.Value(UserContextKey)
	if userVal == nil {
		return nil, exception.NewUnauthorizedException()
	}

	user, ok := userVal.(*entities.User)
	if !ok {
		return nil, exception.NewUnauthorizedException()
	}

	return user, nil
}
