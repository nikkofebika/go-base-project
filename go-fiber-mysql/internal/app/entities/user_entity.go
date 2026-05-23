package entities

import (
	"go-fiber-mysql/internal/app/enums"
	"go-fiber-mysql/internal/app/helpers"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string
	Email       string
	Password    string
	Type        enums.UserType
	LastLoginAt *time.Time

	Roles []Role `gorm:"many2many:user_roles"`

	AuditCreatedEntity
	AuditUpdatedEntity
	AuditDeletedEntity
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.Password != "" {
		password, err := helpers.HashPassword(user.Password)
		if err != nil {
			return err
		}

		user.Password = password
	}

	return nil
}

// func (user *User) MediaID() uint {
// 	return user.ID
// }

// func (user *User) MediaType() enums.MediaType {
// 	return enums.MediaTypeUser
// }
