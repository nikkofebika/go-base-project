package exception

import (
	"go-fiber-mysql/internal/app/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/utils/v2"
)

func ErrorHandler(ctx fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *ValidationException:
		return ctx.Status(e.StatusCode).JSON(helpers.Response[any]{
			Message: e.Message,
			Errors:  e.Errors,
		})
	case *BaseException:
		return ctx.Status(e.StatusCode).JSON(helpers.Response[any]{
			Message: e.Message,
		})
	case *fiber.Error:
		return ctx.Status(e.Code).JSON(helpers.Response[any]{
			Message: e.Error(),
		})
	default:
		statusCode := fiber.StatusInternalServerError
		return ctx.Status(statusCode).JSON(helpers.Response[any]{
			Message: utils.StatusMessage(statusCode),
		})
	}
}
