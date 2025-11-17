package models

import (
	"time"

	"github.com/google/uuid"
)

// Inventory represents an inventory model
type Inventory struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	SKUID     uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_sku_store" json:"sku_id"`
	StoreID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_sku_store" json:"store_id"`
	Quantity  int       `gorm:"not null;default:0" json:"quantity"`
	Version   int       `gorm:"default:1" json:"version"`
	SKU       SKU       `gorm:"foreignKey:SKUID" json:"sku,omitempty"`
	Store     Store     `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Inventory) TableName() string {
	return "inventory"
}
