package requests

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,exists=users.email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenWithPinRequest struct {
	Pin string `json:"pin" validate:"required,exists=users.pin"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
