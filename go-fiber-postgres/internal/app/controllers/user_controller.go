package controllers

import (
	"go-fiber-postgres/internal/app/exception"
	"go-fiber-postgres/internal/app/helpers"
	"go-fiber-postgres/internal/app/helpers/context"
	"go-fiber-postgres/internal/app/requests"
	"go-fiber-postgres/internal/app/responses"
	"go-fiber-postgres/internal/app/services"
	"go-fiber-postgres/internal/pkg/request"
	"go-fiber-postgres/internal/pkg/validator"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	service   services.UserService
	validator *validator.Validator
}

func NewUserController(service services.UserService, validator *validator.Validator) *UserController {
	return &UserController{
		service:   service,
		validator: validator,
	}
}

func (c *UserController) FindAll(ctx fiber.Ctx) error {
	query, err := request.ParseRequestQuery[requests.UserIndexRequest](ctx)
	if err != nil {
		return exception.NewBadRequestException(err.Error())
	}

	req := request.BuildQueryPagination(ctx, query, requests.UserAllowedFilters, requests.UserAllowedIncludes, requests.UserAllowedSorts)

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	datas, meta, err := c.service.FindAll(ctxWithUser, &req)
	if err != nil {
		return err
	}

	return helpers.NewResponsePagination(
		ctx,
		responses.NewUserResponses(datas),
		meta,
	)
}

func (c *UserController) FindOne(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	includes := request.BuildIncludes(ctx, requests.UserAllowedIncludes)

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	data, err := c.service.FindOne(ctxWithUser, id, includes)
	if err != nil {
		return err
	}

	return helpers.NewResponse(ctx, responses.NewUserResponse(&data))
}

func (c *UserController) Create(ctx fiber.Ctx) error {
	body, err := validator.ValidateBody[requests.UserCreateRequest](ctx, c.validator)
	if err != nil {
		return err
	}

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	err = c.service.Create(ctxWithUser, body)
	if err != nil {
		return err
	}

	return helpers.NewResponseCreated(ctx)
}

func (c *UserController) Update(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	body, err := validator.ValidateBody[requests.UserUpdateRequest](ctx, c.validator)
	if err != nil {
		return err
	}

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	err = c.service.Update(ctxWithUser, id, body)
	if err != nil {
		return err
	}

	return helpers.NewResponseUpdated(ctx)
}

func (c *UserController) Delete(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	err = c.service.Delete(ctxWithUser, id)
	if err != nil {
		return err
	}

	return helpers.NewResponseDeleted(ctx)
}

func (c *UserController) ForceDelete(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	err = c.service.ForceDelete(ctxWithUser, id)
	if err != nil {
		return err
	}

	return helpers.NewResponseForceDeleted(ctx)
}

func (c *UserController) Restore(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	err = c.service.Restore(ctxWithUser, id)
	if err != nil {
		return err
	}

	return helpers.NewResponseRestored(ctx)
}

func (c *UserController) SyncRoles(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	body, err := validator.ValidateBody[requests.UserSyncRolesRequest](ctx, c.validator)
	if err != nil {
		return err
	}

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	err = c.service.SyncRoles(ctxWithUser, id, body)
	if err != nil {
		return err
	}

	return helpers.NewResponseUpdated(ctx)
}
