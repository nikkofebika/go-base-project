package enums

type UserType string

const (
	UserTypeAdmin          UserType = "admin"
	UserTypeTechnicianLead UserType = "technician_lead"
	UserTypeTechnician     UserType = "technician"
	UserTypeUser           UserType = "user"
)

func (u UserType) IsValid() bool {
	switch u {
	case UserTypeAdmin, UserTypeTechnicianLead, UserTypeTechnician, UserTypeUser:
		return true
	default:
		return false
	}
}

func (u UserType) GetValues() []string {
	return []string{
		string(UserTypeAdmin),
		string(UserTypeTechnicianLead),
		string(UserTypeTechnician),
		string(UserTypeUser),
	}
}
