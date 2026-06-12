package controllers

import (
	"go-fiber-postgres/internal/app/helpers"
	"go-fiber-postgres/internal/app/helpers/context"
	"go-fiber-postgres/internal/app/requests"
	"go-fiber-postgres/internal/app/responses"
	"go-fiber-postgres/internal/app/services"
	"go-fiber-postgres/internal/pkg/request"
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

func (c *AuthController) ForgotPassword(ctx fiber.Ctx) error {
	body, err := validator.ValidateBody[requests.ForgotPasswordRequest](ctx, c.validator)
	if err != nil {
		return err
	}

	if err = c.service.ForgotPassword(body); err != nil {
		return err
	}

	return helpers.NewResponseMessage(ctx, fiber.StatusOK, "Password reset link has been sent to your email")
}

func (c *AuthController) ResetPassword(ctx fiber.Ctx) error {
	body, err := validator.ValidateBody[requests.ResetPasswordRequest](ctx, c.validator)
	if err != nil {
		return err
	}

	if err = c.service.ResetPassword(body); err != nil {
		return err
	}

	return helpers.NewResponseMessage(ctx, fiber.StatusOK, "Password reset successfully")
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
	request, err := validator.ValidateBody[requests.RefreshTokenRequest](ctx, c.validator)

	if err != nil {
		return err
	}

	token, err := c.service.RefreshToken(request.RefreshToken)
	if err != nil {
		return err
	}

	return helpers.NewResponse(ctx, token)
}

func (c *AuthController) Me(ctx fiber.Ctx) error {
	incudes := request.BuildIncludes(ctx, requests.UserAllowedIncludes)

	context, _, err := context.ContextWithUserID(ctx)
	if err != nil {
		return err
	}

	user, err := c.service.Me(context, incudes)
	if err != nil {
		return err
	}

	return helpers.NewResponse(ctx, responses.NewUserResponse(user))
}
