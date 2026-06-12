package requests

import (
	"go-fiber-postgres/internal/pkg/request"

	"gorm.io/gorm"
)

type PermissionCreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type PermissionUpdateRequest struct {
	Name *string `json:"name" validate:"omitempty,min=2,max=50"`
}

type PermissionIndexRequest struct {
	request.BasePaginationRequest
}

func (r *PermissionIndexRequest) GetBasePaginationRequest() *request.BasePaginationRequest {
	return &r.BasePaginationRequest
}

var PermissionAllowedFilters = request.AllowedFilters{
	"search": func(db *gorm.DB, value string) *gorm.DB {
		return db.Where("name LIKE ?", "%"+value+"%")
	},
}

var PermissionAllowedIncludes = request.AllowedIncludes{}

var PermissionAllowedSorts = request.AllowedSorts{
	"id":         request.SortByColumn("permissions.id"),
	"name":       request.SortByColumn("permissions.name"),
	"created_at": request.SortByColumn("permissions.created_at"),
}
