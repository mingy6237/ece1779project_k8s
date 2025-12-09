package services

import (
	"log"
	"time"

	"inventory-manager-server/config"
	"inventory-manager-server/database"
	"inventory-manager-server/kafka"
	"inventory-manager-server/models"

	"gorm.io/gorm"
)

// OutboxService handles outbox operations
type OutboxService struct {
	// Using global database and kafka instances
}

// NewOutboxService creates a new outbox service
func NewOutboxService() *OutboxService {
	return &OutboxService{}
}

// ProcessOutbox processes outbox records
func (s *OutboxService) ProcessOutbox() error {
	if kafka.Producer == nil {
		// Kafka not initialized, skip processing
		log.Println("Warning: Kafka producer not initialized, skipping outbox processing")
		return nil
	}

	// Query outbox records, filtered by sender_instance_id, ordered by version (maintain consistency)
	var outboxRecords []models.Outbox
	if err := database.DB.Where("sender_instance_id = ?", config.CONFIG.InstanceID).
		Order("version ASC, created_at ASC").
		Limit(100). // Process in batches
		Find(&outboxRecords).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // No records to process
		}
		return err
	}

	if len(outboxRecords) == 0 {
		return nil // No records to process
	}

	// Process each record
	for _, record := range outboxRecords {
		// Send to Kafka
		if err := kafka.PublishMessage(config.CONFIG.KafkaTopic, record); err != nil {
			log.Printf("Failed to publish outbox record %s to Kafka: %v", record.ID, err)
			// Continue processing other records even if one fails
			continue
		}

		// Delete the processed record
		if err := database.DB.Delete(&record).Error; err != nil {
			log.Printf("Failed to delete processed outbox record %s: %v", record.ID, err)
			// Continue processing other records even if deletion fails
			continue
		}

		log.Printf("Processed outbox record %s (SKU: %s, Store: %s, Delta: %d)",
			record.ID, record.SKUName, record.StoreName, record.DeltaQuantity)
	}

	return nil
}

// StartOutboxProcessor starts the outbox processor (background task)
func (s *OutboxService) StartOutboxProcessor() {
	ticker := time.NewTicker(2 * time.Second) // Process every 5 seconds
	defer ticker.Stop()

	log.Println("Outbox processor started")

	// Process immediately on start
	if err := s.ProcessOutbox(); err != nil {
		log.Printf("Error processing outbox on startup: %v", err)
	}

	// Process periodically
	for range ticker.C {
		if err := s.ProcessOutbox(); err != nil {
			log.Printf("Error processing outbox: %v", err)
		}
	}
}
