package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"inventory-manager-server/models"
)

var DB *gorm.DB

// InitDB initializes database connection
func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Only log warnings and errors, not SQL queries
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Store{},
		&models.StoreUser{},
		&models.SKU{},
		&models.Inventory{},
		&models.Outbox{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
