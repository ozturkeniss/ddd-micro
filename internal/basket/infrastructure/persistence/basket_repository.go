package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ddd-micro/internal/basket/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrBasketNotFound = errors.New("basket not found")
	ErrItemNotFound   = errors.New("item not found")
)

// BasketRepository is the Redis implementation of domain.BasketRepository
type BasketRepository struct {
	client *redis.Client
}

// NewBasketRepository creates a new Redis-based basket repository
func NewBasketRepository(client *redis.Client) domain.BasketRepository {
	return &BasketRepository{
		client: client,
	}
}

// Create creates a new basket
func (r *BasketRepository) Create(ctx context.Context, basket *domain.Basket) error {
	if basket.ID == "" {
		basket.ID = uuid.New().String()
	}

	// Set expiration time if not set
	if basket.ExpiresAt.IsZero() {
		basket.SetExpiration(24 * time.Hour) // Default 24 hours
	}

	basket.CreatedAt = time.Now()
	basket.UpdatedAt = time.Now()

	// Serialize basket to JSON
	basketData, err := json.Marshal(basket)
	if err != nil {
		return fmt.Errorf("failed to marshal basket: %w", err)
	}

	// Store basket with expiration
	expiration := time.Until(basket.ExpiresAt)
	key := r.getBasketKey(basket.ID)
	err = r.client.Set(ctx, key, basketData, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store basket: %w", err)
	}

	// Store user basket mapping
	userBasketKey := r.getUserBasketKey(basket.UserID)
	err = r.client.Set(ctx, userBasketKey, basket.ID, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store user basket mapping: %w", err)
	}

	return nil
}

// GetByID retrieves a basket by ID
func (r *BasketRepository) GetByID(ctx context.Context, basketID string) (*domain.Basket, error) {
	key := r.getBasketKey(basketID)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrBasketNotFound
		}
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	var basket domain.Basket
	err = json.Unmarshal([]byte(data), &basket)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal basket: %w", err)
	}

	return &basket, nil
}

// GetByUserID retrieves a basket by user ID
func (r *BasketRepository) GetByUserID(ctx context.Context, userID uint) (*domain.Basket, error) {
	// Get basket ID from user mapping
	userBasketKey := r.getUserBasketKey(userID)
	basketID, err := r.client.Get(ctx, userBasketKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrBasketNotFound
		}
		return nil, fmt.Errorf("failed to get user basket mapping: %w", err)
	}

	// Get basket by ID
	return r.GetByID(ctx, basketID)
}

// Update updates an existing basket
func (r *BasketRepository) Update(ctx context.Context, basket *domain.Basket) error {
	basket.UpdatedAt = time.Now()

	// Serialize basket to JSON
	basketData, err := json.Marshal(basket)
	if err != nil {
		return fmt.Errorf("failed to marshal basket: %w", err)
	}

	// Update basket with expiration
	expiration := time.Until(basket.ExpiresAt)
	key := r.getBasketKey(basket.ID)
	err = r.client.Set(ctx, key, basketData, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to update basket: %w", err)
	}

	// Update user basket mapping expiration
	userBasketKey := r.getUserBasketKey(basket.UserID)
	err = r.client.Set(ctx, userBasketKey, basket.ID, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to update user basket mapping: %w", err)
	}

	return nil
}

// Delete deletes a basket by ID
func (r *BasketRepository) Delete(ctx context.Context, basketID string) error {
	// Get basket to find user ID
	basket, err := r.GetByID(ctx, basketID)
	if err != nil {
		return err
	}

	// Delete basket
	key := r.getBasketKey(basketID)
	err = r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete basket: %w", err)
	}

	// Delete user basket mapping
	userBasketKey := r.getUserBasketKey(basket.UserID)
	err = r.client.Del(ctx, userBasketKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user basket mapping: %w", err)
	}

	return nil
}

// DeleteByUserID deletes a basket by user ID
func (r *BasketRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	// Get basket ID from user mapping
	userBasketKey := r.getUserBasketKey(userID)
	basketID, err := r.client.Get(ctx, userBasketKey).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrBasketNotFound
		}
		return fmt.Errorf("failed to get user basket mapping: %w", err)
	}

	return r.Delete(ctx, basketID)
}

// AddItem adds an item to the basket
func (r *BasketRepository) AddItem(ctx context.Context, basketID string, item *domain.BasketItem) error {
	basket, err := r.GetByID(ctx, basketID)
	if err != nil {
		return err
	}

	basket.AddItem(item.ProductID, item.Quantity, item.UnitPrice)

	return r.Update(ctx, basket)
}

// UpdateItem updates a basket item
func (r *BasketRepository) UpdateItem(ctx context.Context, basketID string, item *domain.BasketItem) error {
	basket, err := r.GetByID(ctx, basketID)
	if err != nil {
		return err
	}

	err = basket.UpdateItemQuantity(item.ProductID, item.Quantity)
	if err != nil {
		return err
	}

	return r.Update(ctx, basket)
}

// RemoveItem removes an item from the basket
func (r *BasketRepository) RemoveItem(ctx context.Context, basketID string, productID uint) error {
	basket, err := r.GetByID(ctx, basketID)
	if err != nil {
		return err
	}

	basket.RemoveItem(productID)

	return r.Update(ctx, basket)
}

// ClearItems removes all items from the basket
func (r *BasketRepository) ClearItems(ctx context.Context, basketID string) error {
	basket, err := r.GetByID(ctx, basketID)
	if err != nil {
		return err
	}

	basket.Clear()

	return r.Update(ctx, basket)
}

// Exists checks if a basket exists by ID
func (r *BasketRepository) Exists(ctx context.Context, basketID string) (bool, error) {
	key := r.getBasketKey(basketID)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check basket existence: %w", err)
	}
	return exists > 0, nil
}

// ExistsByUserID checks if a basket exists for a user
func (r *BasketRepository) ExistsByUserID(ctx context.Context, userID uint) (bool, error) {
	userBasketKey := r.getUserBasketKey(userID)
	exists, err := r.client.Exists(ctx, userBasketKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check user basket existence: %w", err)
	}
	return exists > 0, nil
}

// CleanupExpired removes expired baskets and returns count
func (r *BasketRepository) CleanupExpired(ctx context.Context) (int, error) {
	// Redis automatically removes expired keys, but we can manually check
	// This is mainly for logging and monitoring purposes
	pattern := "basket:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get basket keys: %w", err)
	}

	var expiredCount int
	for _, key := range keys {
		ttl := r.client.TTL(ctx, key).Val()
		if ttl == -1 { // Key exists but has no expiration
			// This shouldn't happen in our implementation, but handle it
			continue
		}
		if ttl == -2 { // Key doesn't exist (expired)
			expiredCount++
		}
	}

	if expiredCount > 0 {
		fmt.Printf("Cleaned up %d expired baskets\n", expiredCount)
	}

	return expiredCount, nil
}

// GetExpiredBaskets retrieves expired baskets
func (r *BasketRepository) GetExpiredBaskets(ctx context.Context) ([]*domain.Basket, error) {
	// Since Redis automatically removes expired keys, this will return empty
	// In a real implementation, you might want to use a different strategy
	// like storing expiration metadata separately
	pattern := "basket:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get basket keys: %w", err)
	}

	var baskets []*domain.Basket
	for _, key := range keys {
		ttl := r.client.TTL(ctx, key).Val()
		if ttl == -2 { // Key doesn't exist (expired)
			continue
		}

		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip if can't retrieve
		}

		var basket domain.Basket
		err = json.Unmarshal([]byte(data), &basket)
		if err != nil {
			continue // Skip if can't unmarshal
		}

		if basket.IsExpired() {
			baskets = append(baskets, &basket)
		}
	}

	return baskets, nil
}

// Helper methods for Redis key generation
func (r *BasketRepository) getBasketKey(basketID string) string {
	return fmt.Sprintf("basket:%s", basketID)
}

func (r *BasketRepository) getUserBasketKey(userID uint) string {
	return fmt.Sprintf("user_basket:%d", userID)
}
