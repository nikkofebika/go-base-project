package jwt

import (
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/enums"
	"go-fiber-mysql/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type JwtTokenDetail struct {
	AccessToken        string    `json:"access_token"`
	AccessTokenExpiry  time.Time `json:"access_token_expiry"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
}

type jwtCustomClaims struct {
	ID    uint           `json:"id"`
	Name  string         `json:"name"`
	Email string         `json:"email"`
	Type  enums.UserType `json:"type"`
	jwt.RegisteredClaims
}

func GenerateToken(user *entities.User, appConfig *config.AppConfig) (JwtTokenDetail, error) {
	jwtTokenDetail := JwtTokenDetail{}

	now := time.Now()
	jwtTokenDetail.AccessTokenExpiry = now.Add(appConfig.JWTExpiresIn)

	accessClaims := jwtCustomClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Type:  user.Type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(jwtTokenDetail.AccessTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	var err error

	jwtTokenDetail.AccessToken, err = accessToken.SignedString([]byte(appConfig.JWTSecret))
	if err != nil {
		return JwtTokenDetail{}, err
	}

	jwtTokenDetail.RefreshTokenExpiry = now.Add(appConfig.JWTRefreshExpiresIn)
	refreshClaims := jwtCustomClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Type:  user.Type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(jwtTokenDetail.RefreshTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	accessToken = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	jwtTokenDetail.RefreshToken, err = accessToken.SignedString([]byte(appConfig.JWTRefreshSecret))
	if err != nil {
		return JwtTokenDetail{}, err
	}

	return jwtTokenDetail, nil
}

func ValidateToken(tokenString, secret string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	})
}

func ExtractUser(token *jwt.Token) (*entities.User, error) {
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	return &entities.User{
		Model: gorm.Model{ID: claims.ID},
		Name:  claims.Name,
		Email: claims.Email,
		Type:  claims.Type,
	}, nil
}

func ExtractUserID(token *jwt.Token) (uint, error) {
	user, err := ExtractUser(token)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
