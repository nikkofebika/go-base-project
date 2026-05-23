package services

import (
	"bytes"
	"context"
	"fmt"
	"go-fiber-mysql/internal/app/entities"
	"go-fiber-mysql/internal/app/enums"
	"go-fiber-mysql/internal/app/exception"
	"go-fiber-mysql/internal/app/helpers"
	"go-fiber-mysql/internal/app/repositories"
	"go-fiber-mysql/internal/app/requests"
	"go-fiber-mysql/internal/config"
	"go-fiber-mysql/internal/pkg/jwt"
	"go-fiber-mysql/internal/pkg/request"
	"html/template"
	"time"
)

type ResetPasswordEmailData struct {
	Name      string
	ResetURL  string
	ExpiresIn string
}

type AuthService interface {
	ForgotPassword(request *requests.ForgotPasswordRequest) error
	Me(userId uint, request *request.PaginationRequest) (*entities.User, error)
	RefreshToken(oldRefreshTokenString string) (*jwt.JwtTokenDetail, error)
	ResetPassword(request *requests.ResetPasswordRequest) error
	Token(request *requests.TokenRequest) (jwt.JwtTokenDetail, error)
}

type authService struct {
	cfg            *config.AppConfig
	emailService   EmailService
	userRepository repositories.UserRepository
	tokenService   TokenService
}

func NewAuthService(cfg *config.AppConfig, emailService EmailService, userRepository repositories.UserRepository, tokenService TokenService) AuthService {
	return &authService{
		cfg:            cfg,
		emailService:   emailService,
		userRepository: userRepository,
		tokenService:   tokenService,
	}
}

func (s *authService) Token(req *requests.TokenRequest) (jwt.JwtTokenDetail, error) {
	user, err := s.userRepository.FindOneByEmail(req.Email)
	if err != nil {
		return jwt.JwtTokenDetail{}, exception.NewBadRequestException("Invalid Credentials!")
	}

	if !helpers.CheckPasswordHash(user.Password, req.Password) {
		return jwt.JwtTokenDetail{}, exception.NewBadRequestException("Invalid Credentials!")
	}

	token, err := jwt.GenerateToken(&user, s.cfg)
	if err != nil {
		return jwt.JwtTokenDetail{}, err
	}

	refreshToken := &entities.Token{
		Type:      enums.TokenTypeRefreshToken,
		UserID:    user.ID,
		Token:     token.RefreshToken,
		ExpiresAt: token.RefreshTokenExpiry,
		IsRevoked: false,
	}

	if err := s.userRepository.SaveRefreshToken(refreshToken); err != nil {
		return jwt.JwtTokenDetail{}, err
	}

	s.userRepository.Update(user.ID, map[string]any{"last_login_at": time.Now()})

	return token, nil
}

func (s *authService) ForgotPassword(request *requests.ForgotPasswordRequest) error {
	user, err := s.userRepository.FindOneByEmail(request.Email)
	if err != nil {
		return exception.NewBadRequestException("Invalid credentials.")
	}

	// Revoke all existing reset tokens for this user (security best practice)
	_ = s.tokenService.RevokeAllUserTokens(user.ID, enums.TokenTypeForgotPassword)

	// Create a new reset token valid for 1 hour
	tokenString, err := s.tokenService.CreateToken(user.ID, time.Hour)
	if err != nil {
		return err
	}

	// Build reset password URL
	resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", s.cfg.WebUrl, tokenString)

	// Render email template
	tmpl, err := template.ParseFiles("internal/templates/emails/reset_password.html")
	if err != nil {
		return exception.NewInternalServerException("Failed to load email template")
	}

	data := ResetPasswordEmailData{
		Name:      user.Name,
		ResetURL:  resetURL,
		ExpiresIn: "1 hour",
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return exception.NewInternalServerException("Failed to render email template")
	}

	// Send email
	subject := "Reset Your Password"
	err = s.emailService.SendEmail([]string{user.Email}, subject, body.String())
	if err != nil {
		return exception.NewInternalServerException("Failed to send reset password email")
	}

	return nil
}

func (s *authService) ResetPassword(request *requests.ResetPasswordRequest) error {
	// Validate token
	token, err := s.tokenService.ValidateToken(request.Token)
	if err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := helpers.HashPassword(request.Password)
	if err != nil {
		return exception.NewInternalServerException("Failed to process password")
	}

	// Update user password
	err = s.userRepository.Update(token.UserID, map[string]any{"password": hashedPassword})
	if err != nil {
		return exception.NewInternalServerException("Failed to update password")
	}

	// Revoke the used token (optional - token will expire anyway, but good for security)
	// _ = s.tokenService.RevokeToken(request.Token)

	// Revoke all other reset tokens for this user (force re-login everywhere)
	_ = s.tokenService.RevokeAllUserTokens(token.UserID, enums.TokenTypeForgotPassword)

	return nil
}

func (s *authService) RefreshToken(oldRefreshTokenString string) (*jwt.JwtTokenDetail, error) {
	// 1. Validasi format token & signature (Stateless check)
	// Ingat: Pakai Refresh Secret, bukan Access Secret
	token, err := jwt.ValidateToken(oldRefreshTokenString, s.cfg.JWTRefreshSecret)
	if err != nil {
		return nil, exception.NewUnauthorizedException()
	}

	// 2. Ambil data user dari token (Claims)
	userClaims, err := jwt.ExtractUser(token)
	if err != nil {
		return nil, exception.NewUnauthorizedException()
	}

	// 3. Cek keberadaan token di Database (Stateful check)
	refreshToken, err := s.userRepository.GetRefreshToken(oldRefreshTokenString)
	if err != nil {
		return nil, exception.NewUnauthorizedException()
	}

	// 4. CEK EXPIRED (Database level)
	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, exception.NewUnauthorizedException()
	}

	// 5. SECURITY: REUSE DETECTION (Rotation Logic)
	// Jika token sudah pernah direvoke tapi dipakai lagi, berarti token dicuri!
	if refreshToken.IsRevoked {
		// Hapus SEMUA refresh token milik user ini (Force Logout Hacker & Victim)

		// Jalankan di background dengan context independen agar tetap berjalan meski request selesai
		go func(userID uint) {
			bgCtx, cancel := context.WithTimeout(context.Background(),
				10*time.Second)
			defer cancel()

			// Pastikan operasi ini tidak terikat dengan context request yang sudah selesai
			_ = bgCtx // bypass context if repo doesn't support it yet, but it's good practice

			if err := s.userRepository.DeleteRefreshToken(userID); err != nil {
				// Gunakan logging untuk memantau kegagalan di background task
				fmt.Printf("[Background Task] Failed to delete all refresh tokens for user %d: %v\n", userID, err)
			}
		}(userClaims.ID)

		return nil, exception.NewUnauthorizedException()
	}

	// 6. ROTASI: Matikan token lama
	refreshToken.IsRevoked = true
	if s.userRepository.SaveRefreshToken(&refreshToken) != nil {
		return nil, exception.NewUnauthorizedException()
	}

	// 7. Ambil data user terbaru (bisi ada perubahan role/status)
	user, err := s.userRepository.FindOne(userClaims.ID, nil)
	if err != nil {
		return nil, exception.NewUnauthorizedException()
	}

	// 8. Generate PAIR Token BARU (Access + Refresh)
	newToken, err := jwt.GenerateToken(&user, s.cfg)
	if err != nil {
		return nil, err
	}

	// 9. Simpan Refresh Token BARU ke Database
	newStoredToken := entities.Token{
		Type:      enums.TokenTypeRefreshToken,
		UserID:    user.ID,
		Token:     newToken.RefreshToken,
		ExpiresAt: newToken.RefreshTokenExpiry,
		IsRevoked: false,
	}

	if err := s.userRepository.SaveRefreshToken(&newStoredToken); err != nil {
		return nil, err
	}

	_ = s.userRepository.Update(user.ID, map[string]any{"last_login_at": time.Now()})

	return &newToken, nil
}

func (s *authService) Me(userId uint, request *request.PaginationRequest) (*entities.User, error) {
	user, err := s.userRepository.FindOne(userId, request.Includes)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
