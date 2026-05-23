package enums

type TokenType string

const (
	TokenTypeRefreshToken   TokenType = "refresh_token"
	TokenTypeForgotPassword TokenType = "technician_lead"
)

func (u TokenType) IsValid() bool {
	switch u {
	case TokenTypeRefreshToken, TokenTypeForgotPassword:
		return true
	default:
		return false
	}
}

func (u TokenType) GetValues() []string {
	return []string{
		string(TokenTypeRefreshToken),
		string(TokenTypeForgotPassword),
	}
}
