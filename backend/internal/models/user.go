package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents user roles
type Role string

const (
	// RoleAdmin is the admin role
	RoleAdmin Role = "admin"
	// RoleUser is the regular user role
	RoleUser Role = "user"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:100;not null;unique" json:"username"`
	Email     string         `gorm:"size:100;not null;unique" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // Password is not included in JSON responses
	Role      Role           `gorm:"size:20;not null;default:'user'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	LastLogin *time.Time     `json:"last_login"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
} 