package services

import (
	"crypto/rand"
	"encoding/hex"
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/enums"
	"go-fiber-mysql/internal/app/exception"
	"go-fiber-mysql/internal/app/repositories"
	"time"

	"gorm.io/gorm"
)

type TokenService interface {
	CreateToken(userID uint64, expiresIn time.Duration) (string, error)
	ValidateToken(tokenString string) (*entities.Token, error)
	RevokeToken(tokenString string) error
	RevokeAllUserTokens(userID uint64, tokenType enums.TokenType) error
}

type tokenService struct {
	repository repositories.TokenRepository
	db         *gorm.DB
}

func NewTokenService(repository repositories.TokenRepository, db *gorm.DB) TokenService {
	return &tokenService{repository, db}
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *tokenService) CreateToken(userID uint64, expiresIn time.Duration) (string, error) {
	tokenString, err := generateSecureToken()
	if err != nil {
		return "", exception.NewInternalServerException("Failed to generate token")
	}

	token := &entities.Token{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(expiresIn),
		IsRevoked: false,
	}

	if err := s.repository.SaveToken(token); err != nil {
		return "", exception.NewInternalServerException("Failed to save token")
	}

	return tokenString, nil
}

func (s *tokenService) ValidateToken(tokenString string) (*entities.Token, error) {
	token, err := s.repository.GetToken(tokenString)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewBadRequestException("Invalid or expired token")
		}
		return nil, exception.NewInternalServerException("Failed to validate token")
	}

	if token.IsRevoked {
		return nil, exception.NewBadRequestException("Token has been revoked")
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, exception.NewBadRequestException("Token has expired")
	}

	return &token, nil
}

func (s *tokenService) RevokeToken(tokenString string) error {
	return s.repository.RevokeToken(tokenString)
}

func (s *tokenService) RevokeAllUserTokens(userID uint64, tokenType enums.TokenType) error {
	return s.repository.DeleteByUserID(userID, tokenType)
}
