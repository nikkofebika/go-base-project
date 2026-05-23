package requests

import (
	"go-fiber-mysql/internal/pkg/request"

	"gorm.io/gorm"
)

type RoleCreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type RoleUpdateRequest struct {
	Name *string `json:"name" validate:"omitempty,min=2,max=50"`
}

type RoleSyncPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" validate:"required,dive,exists=permissions.id"`
}

type RoleIndexRequest struct {
	request.BasePaginationRequest
}

func (r *RoleIndexRequest) GetBasePaginationRequest() *request.BasePaginationRequest {
	return &r.BasePaginationRequest
}

var RoleAllowedFilters = request.AllowedFilters{
	"search": func(db *gorm.DB, value string) *gorm.DB {
		return db.Where("name LIKE ?", "%"+value+"%")
	},
}

var RoleAllowedIncludes = request.AllowedIncludes{
	"permissions": request.IncludePreload("Permissions"),
}

var RoleAllowedSorts = request.AllowedSorts{
	"id":         request.SortByColumn("roles.id"),
	"name":       request.SortByColumn("roles.name"),
	"created_at": request.SortByColumn("roles.created_at"),
}
