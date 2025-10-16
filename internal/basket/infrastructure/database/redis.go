package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     string
	Password string
	DB       string
}

type Database struct {
	Client *redis.Client
}

// NewRedisConnection creates a new Redis database connection
func NewRedisConnection(config Config) (*Database, error) {
	// Parse DB number
	dbNum, err := strconv.Atoi(config.DB)
	if err != nil {
		return nil, fmt.Errorf("invalid redis db number: %w", err)
	}

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       dbNum,
		
		// Connection pool settings
		PoolSize:     10,
		MinIdleConns: 5,
		
		// Timeouts
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		
		// Retry settings
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connection established successfully")

	return &Database{Client: rdb}, nil
}

// Close closes the Redis connection
func (d *Database) Close() error {
	return d.Client.Close()
}

// GetClient returns the Redis client instance
func (d *Database) GetClient() *redis.Client {
	return d.Client
}

// HealthCheck performs a health check on the Redis connection
func (d *Database) HealthCheck(ctx context.Context) error {
	return d.Client.Ping(ctx).Err()
}
