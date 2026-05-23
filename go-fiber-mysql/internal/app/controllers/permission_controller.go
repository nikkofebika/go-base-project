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

type PermissionController struct {
	service   services.PermissionService
	validator *validator.Validator
}

func NewPermissionController(service services.PermissionService, validator *validator.Validator) *PermissionController {
	return &PermissionController{
		service:   service,
		validator: validator,
	}
}

func (c *PermissionController) FindAll(ctx fiber.Ctx) error {
	query, err := request.ParseRequestQuery[requests.PermissionIndexRequest](ctx)
	if err != nil {
		return exception.NewBadRequestException(err.Error())
	}

	req := request.BuildQueryPagination(ctx, query, requests.PermissionAllowedFilters, requests.PermissionAllowedIncludes, requests.PermissionAllowedSorts)

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
		responses.NewPermissionResponses(datas),
		meta,
	)
}

func (c *PermissionController) FindOne(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	includes := request.BuildIncludes(ctx, requests.PermissionAllowedIncludes)

	ctxWithUser, _, err := context.ContextWithUser(ctx)
	if err != nil {
		return err
	}

	data, err := c.service.FindOne(ctxWithUser, id, includes)
	if err != nil {
		return err
	}

	return helpers.NewResponse(ctx, responses.NewPermissionResponse(&data))
}

func (c *PermissionController) Create(ctx fiber.Ctx) error {
	body, err := validator.ValidateBody[requests.PermissionCreateRequest](ctx, c.validator)
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

func (c *PermissionController) Update(ctx fiber.Ctx) error {
	id, err := context.ValidateParamID(ctx, "id")
	if err != nil {
		return err
	}

	body, err := validator.ValidateBody[requests.PermissionUpdateRequest](ctx, c.validator)
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

func (c *PermissionController) Delete(ctx fiber.Ctx) error {
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

func (c *PermissionController) ForceDelete(ctx fiber.Ctx) error {
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

func (c *PermissionController) Restore(ctx fiber.Ctx) error {
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
