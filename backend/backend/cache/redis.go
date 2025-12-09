package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var ctx = context.Background()

// InitRedis initializes Redis connection
func InitRedis(host, port string) error {
	Client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})

	// Test connection
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return nil
}

// Get retrieves a value from cache
func Get(key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Set sets a value in cache
func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// Delete deletes a value from cache
func Delete(key string) error {
	return Client.Del(ctx, key).Err()
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	return Client
}
