package request

import (
	"github.com/gofiber/fiber/v3"
)

const (
	ParamID = "id"
)

type BasePaginationRequest struct {
	Page    int    `query:"page" validate:"min=1"`
	PerPage int    `query:"per_page" validate:"min=1"`
	Sort    string `query:"sort" validate:"omitempty"`
}

type PaginationRequest struct {
	BasePaginationRequest
	Filter   []AppliedFilter
	Includes []AppliedInclude
	Sorts    []AppliedSort
}

type Paginateable interface {
	GetBasePaginationRequest() *BasePaginationRequest
}

func (p *BasePaginationRequest) Normalize() {
	if p.Page <= 1 {
		p.Page = 1
	}

	if p.PerPage <= 1 {
		p.PerPage = 20
	}

	if p.PerPage > 100 {
		p.PerPage = 100
	}
}

func ParseRequestQuery[T any](ctx fiber.Ctx) (*T, error) {
	req := new(T)

	err := ctx.Bind().Query(req)

	return req, err
}

func BuildQueryPagination(ctx fiber.Ctx, query any, allowedFilters AllowedFilters, allowedIncludes AllowedIncludes, allowedSorts AllowedSorts) PaginationRequest {
	var request PaginationRequest

	if allowedFilters != nil {
		request.Filter = ParseFilters(ctx, allowedFilters)
	}

	if allowedIncludes != nil {
		request.Includes = ParseInclude(ctx, allowedIncludes)
	}

	if allowedSorts != nil {
		request.Sorts = ParseSort(ctx, allowedSorts)
	}

	if p, ok := query.(Paginateable); ok {
		pagination := p.GetBasePaginationRequest()
		pagination.Normalize()
		request.BasePaginationRequest = *pagination
	}

	return request
}

func BuildIncludes(ctx fiber.Ctx, allowedIncludes AllowedIncludes) []AppliedInclude {
	return ParseInclude(ctx, allowedIncludes)
}
