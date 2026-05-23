package entities

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name string

	Permissions []Permission `gorm:"many2many:role_permissions"`
	Users       []User       `gorm:"many2many:user_roles"`

	AuditCreatedEntity
	AuditUpdatedEntity
	AuditDeletedEntity
}
