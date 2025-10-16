package domain

// Role represents user roles in the system
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleAdmin:
		return true
	default:
		return false
	}
}

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}

// IsAdmin checks if the role is admin
func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}

// IsUser checks if the role is user
func (r Role) IsUser() bool {
	return r == RoleUser
}
