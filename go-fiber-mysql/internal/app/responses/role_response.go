package responses

import (
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/helpers"
)

type RoleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name,omitempty"`

	Permissions []PermissionResponse `json:"permissions,omitempty"`

	CreatedAt   *string `json:"created_at,omitempty"`
	CreatedByID *uint   `json:"created_by_id,omitempty"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
	UpdatedByID *uint   `json:"updated_by_id,omitempty"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
	DeletedByID *uint   `json:"deleted_by_id,omitempty"`
}

func NewRoleResponse(data *entities.Role) RoleResponse {
	response := RoleResponse{
		ID:   data.ID,
		Name: data.Name,

		CreatedAt:   helpers.ToPointer(entities.TimeToString(data.CreatedAt)),
		CreatedByID: data.CreatedByID,
		UpdatedAt:   helpers.ToPointer(entities.TimeToString(data.UpdatedAt)),
		UpdatedByID: data.UpdatedByID,
		DeletedAt:   entities.DeletedAtToString(data.DeletedAt),
		DeletedByID: data.DeletedByID,
	}

	if len(data.Permissions) > 0 {
		response.Permissions = NewPermissionResponses(data.Permissions)
	}

	return response
}

func NewRoleResponses(datas []entities.Role) []RoleResponse {
	responses := make([]RoleResponse, len(datas))
	for i := range datas {
		responses[i] = NewRoleResponse(&datas[i])
	}
	return responses
}
