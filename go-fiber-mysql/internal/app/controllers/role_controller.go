package controllers

import (
	"go-fiber-mysql/internal/app/exception"
	"go-fiber-mysql/internal/app/helpers"
	"go-fiber-mysql/internal/app/helpers/context"
	"go-fiber-mysql/internal/app/requests"
	"go-fiber-mysql/internal/app/responses"
	"go-fiber-mysql/internal/app/services"
	"go-fiber-mysql/internal/pkg/request"
	"go-fiber-mysql/internal/pkg/validator"

	"github.com/gofiber/fiber/v3"
)

type RoleController struct {
	service   services.RoleService
	validator *validator.Validator
}

func NewRoleController(service services.RoleService, validator *validator.Validator) *RoleController {
	return &RoleController{
		service:   service,
		validator: validator,
	}
}

func (c *RoleController) FindAll(ctx fiber.Ctx) error {
	query, err := request.ParseRequestQuery[requests.RoleIndexRequest](ctx)
	if err != nil {
		return exception.NewBadRequestException(err.Error())
	}

	req := request.BuildQueryPagination(ctx, query, requests.RoleAllowedFilters, requests.RoleAllowedIncludes, requests.RoleAllowedSorts)

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
		responses.NewRoleResponses(datas),
		meta,
	)
}

func (c *RoleController) FindOne(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	includes := request.BuildIncludes(ctx, requests.RoleAllowedIncludes)

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	data, err := c.service.FindOne(ctxWithUser, id, includes)
	if err != nil {
		return err
	}

	return helpers.NewResponse(ctx, responses.NewRoleResponse(&data))
}

func (c *RoleController) Create(ctx fiber.Ctx) error {
	body, err := validator.ValidateBody[requests.RoleCreateRequest](ctx, c.validator)
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

func (c *RoleController) Update(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	body, err := validator.ValidateBody[requests.RoleUpdateRequest](ctx, c.validator)
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

func (c *RoleController) Delete(ctx fiber.Ctx) error {
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

func (c *RoleController) ForceDelete(ctx fiber.Ctx) error {
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

	return helpers.NewResponseRestored(ctx)
}

func (c *RoleController) Restore(ctx fiber.Ctx) error {
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

func (c *RoleController) SyncPermissions(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	body, err := validator.ValidateBody[requests.RoleSyncPermissionsRequest](ctx, c.validator)
	if err != nil {
		return err
	}

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	err = c.service.SyncPermissions(ctxWithUser, id, body)
	if err != nil {
		return err
	}

	return helpers.NewResponseUpdated(ctx)
}
