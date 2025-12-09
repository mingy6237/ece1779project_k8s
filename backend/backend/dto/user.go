package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse represents a user in API responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserPreviewResponse represents a user in API responses
type UserPreviewResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=manager staff"`
}

// UpdateUserRequest represents a request to update user information
type UpdateUserRequest struct {
	TargetID uuid.UUID `json:"target_id" binding:"required"`
	Username string    `json:"username" binding:"omitempty"`
	Email    string    `json:"email" binding:"omitempty,email"`
	Role     string    `json:"role" binding:"omitempty,oneof=manager staff"`
}

// ChangePasswordRequest represents a request to change password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
