package models

import (
	"time"

	"github.com/google/uuid"
)

// StoreUser represents a store-user association model
type StoreUser struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null;index;uniqueIndex:idx_store_user" json:"user_id"`
	StoreID   uuid.UUID `gorm:"column:store_id;type:uuid;not null;index;uniqueIndex:idx_store_user" json:"store_id"`
	Version   int       `gorm:"default:1" json:"version"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Store     Store     `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (StoreUser) TableName() string {
	return "store_user"
}
