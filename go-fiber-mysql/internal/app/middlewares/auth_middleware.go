package middlewares

import (
	"go-fiber-mysql/internal/app/exception"
	"go-fiber-mysql/internal/config"
	"go-fiber-mysql/internal/pkg/jwt"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(cfg *config.AppConfig) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		// get authorization bearer
		header := ctx.Get(fiber.HeaderAuthorization)
		if header == "" {
			return exception.NewUnauthorizedException()
		}

		// split and take the bearer token
		tokens := strings.Split(header, " ")
		if len(tokens) != 2 || tokens[0] != "Bearer" {
			return exception.NewUnauthorizedException()
		}

		// validate token
		jwtToken, err := jwt.ValidateToken(tokens[1], cfg.JWTSecret)
		if err != nil {
			return exception.NewUnauthorizedException()
		}

		// get data from jwtToken
		user, err := jwt.ExtractUser(jwtToken)
		if err != nil {
			return exception.NewUnauthorizedException()
		}

		// attach data to the context
		ctx.Locals("user", user)
		ctx.Locals("user_id", user.ID)
		return ctx.Next()
	}
}
