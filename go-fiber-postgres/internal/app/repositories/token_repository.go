package repositories

import (
	"go-fiber-postgres/internal/app/entities"
	"go-fiber-postgres/internal/app/enums"

	"gorm.io/gorm"
)

type TokenRepository interface {
	WithTx(db *gorm.DB) TokenRepository
	DB() *gorm.DB
	GetToken(token string) (entities.Token, error)
	GetRefreshToken(token string) (entities.Token, error)
	DeleteByUserID(userID uint64, tokenType enums.TokenType) error
	DeleteByToken(token string) error
	DeleteRefreshToken(token string) error
	DeleteAllByUserID(userID uint64) error
	RevokeToken(token string) error
	SaveToken(token *entities.Token) error
	SaveRefreshToken(token *entities.Token) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) WithTx(db *gorm.DB) TokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) DB() *gorm.DB {
	return r.db.Model(&entities.Token{})
}

func (r *tokenRepository) GetToken(token string) (entities.Token, error) {
	var data entities.Token
	err := r.db.Where("token = ?", token).First(&data).Error
	return data, err
}

func (r *tokenRepository) GetRefreshToken(token string) (entities.Token, error) {
	var data entities.Token
	err := r.db.Where("token = ?", token).Where("type = ?", enums.TokenTypeRefreshToken).First(&data).Error
	return data, err
}

func (r *tokenRepository) SaveToken(token *entities.Token) error {
	return r.db.Save(token).Error
}

func (r *tokenRepository) SaveRefreshToken(token *entities.Token) error {
	return r.db.Save(token).Error
}

func (r *tokenRepository) DeleteByUserID(userID uint64, tokenType enums.TokenType) error {
	return r.db.Where("user_id = ?", userID).Where("type=?", tokenType).Delete(&entities.Token{}).Error
}

func (r *tokenRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&entities.Token{}).Error
}

func (r *tokenRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Where("type = ?", enums.TokenTypeRefreshToken).Delete(&entities.Token{}).Error
}

func (r *tokenRepository) DeleteAllByUserID(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&entities.Token{}).Error
}

func (r *tokenRepository) RevokeToken(token string) error {
	return r.db.Model(&entities.Token{}).Where("token = ?", token).Update("is_revoked", true).Error
}
