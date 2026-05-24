package entities

type Role struct {
	IDEntity
	Name string

	Permissions []Permission `gorm:"many2many:role_permissions"`
	Users       []User       `gorm:"many2many:user_roles"`

	AuditCreatedEntity
	AuditUpdatedEntity
	AuditDeletedEntity
}
