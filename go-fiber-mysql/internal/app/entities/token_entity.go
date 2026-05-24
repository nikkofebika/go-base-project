package entities

import (
	"go-fiber-mysql/internal/app/enums"
	"time"
)

type Token struct {
	IDEntity
	UserID    uint64
	Type      enums.TokenType
	Token     string
	ExpiresAt time.Time
	IsRevoked bool // rotation here

	User User `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
