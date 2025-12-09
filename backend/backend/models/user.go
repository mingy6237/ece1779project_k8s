package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user model
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null;size:50" json:"username"`
	PasswordHash string    `gorm:"not null;size:255" json:"-"`
	Email        string    `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Role         string    `gorm:"not null;size:20" json:"role"` // 'manager' or 'staff'
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
