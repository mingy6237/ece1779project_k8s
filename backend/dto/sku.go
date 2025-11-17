package dto

import (
	"time"

	"github.com/google/uuid"
)

// SKUResponse represents a SKU in API responses
type SKUResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSKURequest represents a request to create a new SKU
type CreateSKURequest struct {
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
}

// UpdateSKURequest represents a request to update SKU information
type UpdateSKURequest struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,min=0"`
}
