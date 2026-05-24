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

type IDEntity struct {
	ID uint64 `gorm:"primarykey" json:"id"`
}

type CreatedAtEntity struct {
	CreatedAt time.Time `json:"created_at"`
}

type CreatedByIDEntity struct {
	CreatedByID *uint64 `json:"created_by_id"`
}

type AuditCreatedEntity struct {
	CreatedAtEntity
	CreatedByIDEntity
}

type UpdatedAtEntity struct {
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdatedByIDEntity struct {
	UpdatedByID *uint64 `json:"updated_by_id"`
}

type AuditUpdatedEntity struct {
	UpdatedAtEntity
	UpdatedByIDEntity
}

type DeletedAtEntity struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type DeletedByIDEntity struct {
	DeletedByID *uint64 `json:"deleted_by_id"`
}

type AuditDeletedEntity struct {
	DeletedAtEntity
	DeletedByIDEntity
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
