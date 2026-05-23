package requests

import (
	"go-fiber-mysql/internal/pkg/request"

	"gorm.io/gorm"
)

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=250"`
	Email    string `json:"email" validate:"required,email,min=2,max=250"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

type UserUpdateRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=2,max=255"`
	Email    *string `json:"email" validate:"omitempty,email,min=2,max=255"`
	Password *string `json:"password" validate:"omitempty,min=6,max=255"`
}

type UserSyncRolesRequest struct {
	RoleIDs []uint `json:"role_ids" validate:"required,dive,exists=roles.id"`
}

// pagination request
type UserIndexRequest struct {
	request.BasePaginationRequest
}

func (r *UserIndexRequest) GetBasePaginationRequest() *request.BasePaginationRequest {
	return &r.BasePaginationRequest
}

var UserAllowedFilters = request.AllowedFilters{
	"search": func(db *gorm.DB, value string) *gorm.DB {
		return db.Where("name LIKE ? OR email LIKE ?", "%"+value+"%", "%"+value+"%")
	},
}

var UserAllowedIncludes = request.AllowedIncludes{
	"roles":             request.IncludePreload("Roles"),
	"roles.permissions": request.IncludePreload("Roles.Permissions"),
}

var UserAllowedSorts = request.AllowedSorts{
	"id":         request.SortByColumn("users.id"),
	"created_at": request.SortByColumn("users.created_at"),
	"updated_at": request.SortByColumn("users.updated_at"),
}
