package consumers

import (
	"context"
	"fmt"
	"log"

	"github.com/ddd-micro/kafka"
)

// BasketConsumer handles basket-related Kafka events
type BasketConsumer struct {
	eventPublisher kafka.EventPublisher
}

// NewBasketConsumer creates a new basket consumer
func NewBasketConsumer(eventPublisher kafka.EventPublisher) *BasketConsumer {
	return &BasketConsumer{
		eventPublisher: eventPublisher,
	}
}

// HandlePaymentCompleted handles payment completed events
func (c *BasketConsumer) HandlePaymentCompleted(ctx context.Context, event kafka.PaymentCompletedEvent) error {
	log.Printf("Processing payment completed event: %s", event.Data.PaymentID)

	// Clear basket if it was a basket-based payment
	if event.Data.BasketID != nil {
		// Convert payment items to basket items for the cleared event
		var basketItems []kafka.PaymentItem
		for _, item := range event.Data.Items {
			basketItems = append(basketItems, kafka.PaymentItem{
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				UnitPrice:  item.UnitPrice,
				TotalPrice: item.TotalPrice,
			})
		}

		// Publish basket cleared event
		basketClearedEvent := kafka.BasketClearedEvent{
			BaseEvent: kafka.NewBaseEvent(kafka.EventTypeBasketCleared, "basket-service"),
			Data: kafka.BasketClearedData{
				UserID:    event.Data.UserID,
				BasketID:  *event.Data.BasketID,
				Items:     basketItems,
				Reason:    "payment_completed",
				OrderID:   &event.Data.OrderID,
				PaymentID: &event.Data.PaymentID,
			},
		}

		if err := c.eventPublisher.PublishBasketCleared(basketClearedEvent); err != nil {
			return fmt.Errorf("failed to publish basket cleared event: %w", err)
		}

		log.Printf("Basket cleared for user %d, basket %s", event.Data.UserID, *event.Data.BasketID)
	}

	return nil
}

// HandlePaymentFailed handles payment failed events
func (c *BasketConsumer) HandlePaymentFailed(ctx context.Context, event kafka.PaymentFailedEvent) error {
	log.Printf("Processing payment failed event: %s", event.Data.PaymentID)

	// Release basket reservation if it was a basket-based payment
	if event.Data.BasketID != nil {
		log.Printf("Payment failed for basket %s, reservation should be released", *event.Data.BasketID)
		// Here you would implement basket reservation release logic
		// This might involve calling the basket service API or publishing another event
	}

	return nil
}

// HandlePaymentCancelled handles payment cancelled events
func (c *BasketConsumer) HandlePaymentCancelled(ctx context.Context, event kafka.PaymentCancelledEvent) error {
	log.Printf("Processing payment cancelled event: %s", event.Data.PaymentID)

	// Release basket reservation if it was a basket-based payment
	if event.Data.BasketID != nil {
		log.Printf("Payment cancelled for basket %s, reservation should be released", *event.Data.BasketID)
		// Here you would implement basket reservation release logic
	}

	return nil
}

// Start starts the basket consumer
func (c *BasketConsumer) Start(ctx context.Context) error {
	log.Println("Starting basket consumer...")

	// In a real implementation, you would:
	// 1. Create a Kafka consumer
	// 2. Subscribe to payment events
	// 3. Handle events in a loop

	// For now, this is a placeholder
	log.Println("Basket consumer started successfully")
	return nil
}

// Stop stops the basket consumer
func (c *BasketConsumer) Stop() error {
	log.Println("Stopping basket consumer...")
	return nil
}
