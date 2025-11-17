package models

import (
	"time"

	"github.com/google/uuid"
)

// Store represents a store model
type Store struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"not null;size:100" json:"name"`
	Address   string    `gorm:"not null;size:255" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Store) TableName() string {
	return "stores"
}
