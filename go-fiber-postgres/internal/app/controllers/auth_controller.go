package controllers

import (
	"go-fiber-postgres/internal/app/helpers"
	"go-fiber-postgres/internal/app/requests"
	"go-fiber-postgres/internal/app/services"
	"go-fiber-postgres/internal/pkg/validator"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	service   services.AuthService
	validator *validator.Validator
}

func NewAuthController(service services.AuthService, validator *validator.Validator) *AuthController {
	return &AuthController{
		service:   service,
		validator: validator,
	}
}

func (c *AuthController) Token(ctx fiber.Ctx) error {
	body, err := validator.ValidateBody[requests.TokenRequest](ctx, c.validator)
	if err != nil {
		return err
	}

	data, err := c.service.Token(body)
	if err != nil {
		return err
	}

	return helpers.NewResponse(ctx, data)
}

func (c *AuthController) RefreshToken(ctx fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"data": "RefreshToken",
	})
}
