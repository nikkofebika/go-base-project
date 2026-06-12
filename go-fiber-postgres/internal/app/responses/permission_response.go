package responses

import (
	"go-fiber-postgres/internal/app/entities"
	"go-fiber-postgres/internal/app/helpers"
)

type PermissionResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name,omitempty"`

	CreatedAt   *string `json:"created_at,omitempty"`
	CreatedByID *uint64 `json:"created_by_id,omitempty"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
	UpdatedByID *uint64 `json:"updated_by_id,omitempty"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
	DeletedByID *uint64 `json:"deleted_by_id,omitempty"`
}

func NewPermissionResponse(data *entities.Permission) PermissionResponse {
	return PermissionResponse{
		ID:   data.ID,
		Name: data.Name,

		CreatedAt:   helpers.ToPointer(entities.TimeToString(data.CreatedAt)),
		CreatedByID: data.CreatedByID,
		UpdatedAt:   helpers.ToPointer(entities.TimeToString(data.UpdatedAt)),
		UpdatedByID: data.UpdatedByID,
		DeletedAt:   entities.DeletedAtToString(data.DeletedAt),
		DeletedByID: data.DeletedByID,
	}
}

func NewPermissionResponses(datas []entities.Permission) []PermissionResponse {
	responses := make([]PermissionResponse, len(datas))
	for i := range datas {
		responses[i] = NewPermissionResponse(&datas[i])
	}
	return responses
}
