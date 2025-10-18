package kafka

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds Kafka configuration
type Config struct {
	Brokers       []string
	Topic         string
	GroupID       string
	Offset        int64
	RetryAttempts int
	RetryDelay    time.Duration
	Timeout       time.Duration
}

// LoadConfig loads Kafka configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Brokers:       getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		Topic:         getEnv("KAFKA_TOPIC", "payment-events"),
		GroupID:       getEnv("KAFKA_GROUP_ID", "payment-service"),
		Offset:        getEnvAsInt64("KAFKA_OFFSET", sarama.OffsetNewest),
		RetryAttempts: getEnvAsInt("KAFKA_RETRY_ATTEMPTS", 3),
		RetryDelay:    getEnvAsDuration("KAFKA_RETRY_DELAY", 5*time.Second),
		Timeout:       getEnvAsDuration("KAFKA_TIMEOUT", 30*time.Second),
	}
}

// GetPublisherConfig returns publisher configuration
func (c *Config) GetPublisherConfig() *PublisherConfig {
	return &PublisherConfig{
		Brokers: c.Brokers,
		Topic:   c.Topic,
		Timeout: c.Timeout,
	}
}

// GetConsumerConfig returns consumer configuration
func (c *Config) GetConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		Brokers:       c.Brokers,
		Topic:         c.Topic,
		GroupID:       c.GroupID,
		Offset:        c.Offset,
		RetryAttempts: c.RetryAttempts,
		RetryDelay:    c.RetryDelay,
	}
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
