package dto

import (
	"time"

	"github.com/google/uuid"
)

// StoreResponse represents a store in API responses
type StoreResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateStoreRequest represents a request to create a new store
type CreateStoreRequest struct {
	Name    string `json:"name" binding:"required,max=100"`
	Address string `json:"address" binding:"required,max=255"`
}

// UpdateStoreRequest represents a request to update store information
type UpdateStoreRequest struct {
	Name    string `json:"name" binding:"omitempty,max=100"`
	Address string `json:"address" binding:"omitempty,max=255"`
}

// AddStaffToStoreRequest represents a request to add staff to a store
type AddStaffToStoreRequest struct {
	StoreID uuid.UUID `json:"store_id" binding:"required"`
	UserID  uuid.UUID `json:"user_id" binding:"required"`
}

// StoreStaffResponse represents a store staff member in API responses
type StoreStaffResponse struct {
	ID        uuid.UUID     `json:"id"`
	StoreID   uuid.UUID     `json:"store_id"`
	UserID    uuid.UUID     `json:"user_id"`
	User      UserResponse  `json:"user"`
	Store     StoreResponse `json:"store"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// StoreStaffListResponse represents the response for listing store staff
type StoreStaffListResponse struct {
	StoreID uuid.UUID      `json:"store_id"`
	Staff   []UserResponse `json:"staff"`
}
