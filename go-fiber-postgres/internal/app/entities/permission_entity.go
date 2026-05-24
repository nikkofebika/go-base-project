package entities

type Permission struct {
	IDEntity
	Name string
	Slug string `gorm:"uniqueIndex"`

	AuditCreatedEntity
	AuditUpdatedEntity
	AuditDeletedEntity
}
