package entities

import (
	"go-fiber-postgres/internal/app/enums"
	"time"
)

type Token struct {
	ID        uint
	UserID    uint64
	Type      enums.TokenType
	Token     string
	ExpiresAt time.Time
	IsRevoked bool // rotation here

	User User `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
