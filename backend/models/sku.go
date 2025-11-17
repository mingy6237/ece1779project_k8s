package models

import (
	"time"

	"github.com/google/uuid"
)

// SKU represents a product model
type SKU struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"not null;size:100" json:"name"`
	Category    string    `gorm:"size:100" json:"category"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(12,2)" json:"price"`
	Version     int       `gorm:"default:1" json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SKU) TableName() string {
	return "sku"
}
