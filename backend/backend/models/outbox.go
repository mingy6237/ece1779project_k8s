package models

import (
	"time"

	"github.com/google/uuid"
)

// Outbox represents an outbox record model (for transactional outbox pattern)
type Outbox struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OperationType    string    `gorm:"not null;size:20" json:"operation_type"` // "create", "update", "adjust", "delete"
	SenderInstanceID string    `gorm:"not null;size:100;index" json:"sender_instance_id"`
	InventoryID      uuid.UUID `gorm:"column:inventory_id;type:uuid;not null;index" json:"inventory_id"`
	SKUID            uuid.UUID `gorm:"column:sku_id;type:uuid;not null;index" json:"sku_id"`
	SKUName          string    `gorm:"not null;size:100" json:"sku_name"`
	StoreID          uuid.UUID `gorm:"column:store_id;type:uuid;not null;index" json:"store_id"`
	StoreName        string    `gorm:"not null;size:100" json:"store_name"`
	UserID           uuid.UUID `gorm:"column:user_id;type:uuid;not null" json:"user_id"`
	UserName         string    `gorm:"not null;size:100" json:"user_name"`
	DeltaQuantity    int       `gorm:"not null" json:"delta_quantity"`
	NewQuantity      int       `gorm:"not null" json:"new_quantity"`
	Version          int       `gorm:"default:1" json:"version"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (Outbox) TableName() string {
	return "outbox"
}
