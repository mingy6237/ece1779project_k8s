package config

import (
	"os"
)

type Config struct {
	// Instance info
	InstanceID string

	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis configuration
	RedisHost string
	RedisPort string

	// Kafka configuration
	KafkaBrokers string
	KafkaTopic   string

	// Server configuration
	ServerPort string
}

var CONFIG *Config

func LoadConfig() *Config {
	instanceID := getEnv("INSTANCE_ID", "")
	if instanceID == "" {
		panic("INSTANCE_ID is not set")
	}
	CONFIG = &Config{
		InstanceID: instanceID,

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "inventory_db"),

		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),

		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "inventory-updates"),

		ServerPort: getEnv("SERVER_PORT", "3000"),
	}
	return CONFIG
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
