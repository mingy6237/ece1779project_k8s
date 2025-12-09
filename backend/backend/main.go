package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"

	"inventory-manager-server/cache"
	"inventory-manager-server/config"
	"inventory-manager-server/database"
	"inventory-manager-server/events"
	"inventory-manager-server/kafka"
	"inventory-manager-server/models"
	"inventory-manager-server/routes"
	"inventory-manager-server/services"
	"inventory-manager-server/utils"
	"inventory-manager-server/websocket"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Initialize configuration
	cfg := config.LoadConfig()
	log.Printf("Configuration loaded: Server Port=%s, DB=%s, Redis=%s:%s",
		cfg.ServerPort, cfg.DBName, cfg.RedisHost, cfg.RedisPort)

	// Initialize JWT
	utils.InitJWT()

	// GORM PostgreSQL DSN format: postgres://user:password@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	if err := database.InitDB(dsn); err != nil {
		log.Printf("Warning: Failed to initialize database: %v (continuing without database)", err)
	}

	// Initialize Redis connection
	if err := cache.InitRedis(cfg.RedisHost, cfg.RedisPort); err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v (continuing without cache)", err)
	}

	// Initialize Kafka producer/consumer
	if err := kafka.InitProducer([]string{cfg.KafkaBrokers}); err != nil {
		log.Printf("Warning: Failed to initialize Kafka producer: %v (continuing without Kafka producer)", err)
	}
	if err := kafka.InitConsumer([]string{cfg.KafkaBrokers}); err != nil {
		log.Printf("Warning: Failed to initialize Kafka consumer: %v (continuing without Kafka consumer)", err)
	}
	if err := kafka.StartConsumer(context.Background(), cfg.KafkaTopic, events.InventoryUpdateHandler); err != nil {
		log.Printf("Warning: Failed to start Kafka consumer: %v (continuing without Kafka consumer)", err)
	}

	// Initialize WebSocket Hub
	if err := websocket.InitHub(); err != nil {
		log.Printf("Warning: Failed to initialize WebSocket Hub: %v (continuing without WebSocket)", err)
	}

	// Start background tasks
	// Start Outbox processor
	outboxService := services.NewOutboxService()
	go outboxService.StartOutboxProcessor()
	log.Println("Outbox processor background task started")

	// Create admin user if not exists
	var adminUser models.User
	result := database.DB.Where("username = ?", "admin").First(&adminUser)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("adminadmin"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Warning: Failed to hash password: %v (continuing without admin user)", err)
		} else {
			user := models.User{
				Username:     "admin",
				PasswordHash: string(hashedPassword),
				Email:        "admin@admin.com",
				Role:         "manager",
			}
			if err := database.DB.Create(&user).Error; err != nil {
				log.Printf("Warning: Failed to create admin user: %v (continuing without admin user)", err)
			} else {
				log.Println("Admin user created.")
			}
		}
	} else if result.Error == nil {
		log.Println("Admin user already exists. Skipping creation.")
	} else {
		log.Printf("Warning: Could not verify if admin user exists: %v", result.Error)
	}
	// Setup routes
	router := routes.SetupRoutes()

	// Start HTTP server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
