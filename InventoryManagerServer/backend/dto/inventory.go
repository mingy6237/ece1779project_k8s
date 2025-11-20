package dto

import (
	"time"

	"github.com/google/uuid"
)

// InventoryResponse represents inventory in API responses
type InventoryResponse struct {
	ID        uuid.UUID      `json:"id"`
	SKUID     uuid.UUID      `json:"sku_id"`
	StoreID   uuid.UUID      `json:"store_id"`
	Quantity  int            `json:"quantity"`
	Version   int            `json:"version"`
	SKU       *SKUResponse   `json:"sku,omitempty"`
	Store     *StoreResponse `json:"store,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// CreateInventoryRequest represents a request to create new inventory
type CreateInventoryRequest struct {
	SKUID    uuid.UUID `json:"sku_id" binding:"required"`
	StoreID  uuid.UUID `json:"store_id" binding:"required"`
	Quantity int       `json:"quantity" binding:"required,min=0"`
}

// UpdateInventoryRequest represents a request to update inventory
type UpdateInventoryRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

// InventoryQueryParams represents query parameters for inventory listing
type InventoryQueryParams struct {
	StoreID  *string `form:"store_id"` // Single store ID from query parameter, nil means all stores
	SKUID    *string `form:"sku_id"`
	Page     int     `form:"page,default=1" binding:"min=1"`
	PageSize int     `form:"page_size,default=20" binding:"min=1,max=100"`
	SortBy   string  `form:"sort_by,default=created_at"`
	Order    string  `form:"order,default=desc" binding:"oneof=asc desc"`
}

// InventoryListResponse represents the response for listing inventory
type InventoryListResponse struct {
	Items      []InventoryResponse `json:"items"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

// AdjustInventoryRequest represents a request to adjust inventory quantity
type AdjustInventoryRequest struct {
	DeltaQuantity int `json:"delta_quantity" binding:"required"`
}
