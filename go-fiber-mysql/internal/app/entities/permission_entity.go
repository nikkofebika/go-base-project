package entities

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name string
	Slug string `gorm:"uniqueIndex"`

	AuditCreatedEntity
	AuditUpdatedEntity
	AuditDeletedEntity
}
