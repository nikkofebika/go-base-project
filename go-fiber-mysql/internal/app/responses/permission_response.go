package responses

import (
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/helpers"
)

type PermissionResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`

	CreatedAt   *string `json:"created_at,omitempty"`
	CreatedByID *uint   `json:"created_by_id,omitempty"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
	UpdatedByID *uint   `json:"updated_by_id,omitempty"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
	DeletedByID *uint   `json:"deleted_by_id,omitempty"`
}

func NewPermissionResponse(data *entities.Permission) PermissionResponse {
	return PermissionResponse{
		ID:   data.ID,
		Name: data.Name,
		Slug: data.Slug,

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
