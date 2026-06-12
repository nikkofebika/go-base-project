package entities

type Permission struct {
	IDEntity
	Name string `gorm:"uniqueIndex"`

	AuditCreatedEntity
	AuditUpdatedEntity
	AuditDeletedEntity
}
