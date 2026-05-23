package responses

import (
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/helpers"
	"time"
)

type UserResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name,omitempty"`
	Email       string     `json:"email,omitempty"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`

	Roles []RoleResponse `json:"roles,omitempty"`

	CreatedAt   *string `json:"created_at,omitempty"`
	CreatedByID *uint   `json:"created_by_id,omitempty"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
	UpdatedByID *uint   `json:"updated_by_id,omitempty"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
	DeletedByID *uint   `json:"deleted_by_id,omitempty"`
}

func NewUserResponse(data *entities.User) UserResponse {
	response := UserResponse{
		ID:          data.ID,
		Name:        data.Name,
		Email:       data.Email,
		LastLoginAt: data.LastLoginAt,

		CreatedAt:   helpers.ToPointer(entities.TimeToString(data.CreatedAt)),
		CreatedByID: data.CreatedByID,
		UpdatedAt:   helpers.ToPointer(entities.TimeToString(data.UpdatedAt)),
		UpdatedByID: data.UpdatedByID,
		DeletedAt:   entities.DeletedAtToString(data.DeletedAt),
		DeletedByID: data.DeletedByID,
	}

	if len(data.Roles) > 0 {
		response.Roles = NewRoleResponses(data.Roles)
	}

	return response
}

func NewUserResponses(datas []entities.User) []UserResponse {
	responses := make([]UserResponse, len(datas))
	for i := range datas {
		responses[i] = NewUserResponse(&datas[i])
	}
	return responses
}
