package entities

import (
	"time"

	"gorm.io/gorm"
)

const (
	CreatedAt   = "created_at"
	CreatedByID = "created_by_id"
	UpdatedAt   = "updated_at"
	UpdatedByID = "updated_by_id"
	DeletedAt   = "deleted_at"
	DeletedByID = "deleted_by_id"
)

type UpdatedFields map[string]any

type AuditCreatedEntity struct {
	CreatedByID *uint `json:"created_by_id"`
}

type AuditUpdatedEntity struct {
	UpdatedByID *uint `json:"updated_by_id"`
}

type AuditDeletedEntity struct {
	DeletedByID *uint `json:"deleted_by_id"`
}

func TimeToString(t time.Time) string {
	return t.Format(time.DateTime)
}

func DeletedAtToString(t gorm.DeletedAt) *string {
	if !t.Valid {
		return nil
	}

	value := t.Time.Format(time.DateTime)
	return &value
}
