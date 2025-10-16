package domain

import (
	"time"

	"gorm.io/gorm"
)

// Test comment with bad formatting

// User represents the user domain entity
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string         `gorm:"not null;size:255" json:"-"`
	FirstName string         `gorm:"size:100" json:"first_name"`
	LastName  string         `gorm:"size:100" json:"last_name"`
	Role      Role           `gorm:"type:varchar(20);default:'user'" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for User entity
func (User) TableName() string {
	return "users"
}

// IsValidEmail checks if the email is valid
func (u *User) IsValidEmail() bool {
	return len(u.Email) > 0 && len(u.Email) <= 255
}

// IsValidPassword checks if the password meets minimum requirements
func (u *User) IsValidPassword() bool {
	return len(u.Password) >= 6
}

// Activate activates the user account
func (u *User) Activate() {
	u.IsActive = true
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role.IsAdmin()
}

// HasRole checks if the user has the specified role
func (u *User) HasRole(role Role) bool {
	return u.Role == role
}

// AssignRole assigns a new role to the user
func (u *User) AssignRole(role Role) {
	if role.IsValid() {
		u.Role = role
	}
}
