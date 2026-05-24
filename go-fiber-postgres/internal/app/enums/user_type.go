package enums

type UserType string

const (
	UserTypeAdmin UserType = "admin" // given all access, without permissions or roles
	UserTypeUser  UserType = "user"  // must have permissions and roles, to access the system
)

func (u UserType) IsValid() bool {
	switch u {
	case UserTypeAdmin, UserTypeUser:
		return true
	default:
		return false
	}
}

func (u UserType) GetValues() []string {
	return []string{
		string(UserTypeAdmin),
		string(UserTypeUser),
	}
}
