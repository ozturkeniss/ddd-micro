package consumers

import (
	"context"
	"log"

	"github.com/ddd-micro/kafka"
)

// ProductConsumer handles product-related Kafka events
type ProductConsumer struct {
	eventPublisher kafka.EventPublisher
}

// NewProductConsumer creates a new product consumer
func NewProductConsumer(eventPublisher kafka.EventPublisher) *ProductConsumer {
	return &ProductConsumer{
		eventPublisher: eventPublisher,
	}
}

// HandlePaymentCompleted handles payment completed events
func (c *ProductConsumer) HandlePaymentCompleted(ctx context.Context, event kafka.PaymentCompletedEvent) error {
	log.Printf("Processing payment completed event for stock update: %s", event.Data.PaymentID)

	// Update stock for each item in the payment
	for _, item := range event.Data.Items {
		// Calculate new stock (assuming we have the current stock)
		// In a real implementation, you would fetch current stock from the product service
		newStock := 100 - item.Quantity // Placeholder calculation

		// Publish stock updated event
		stockUpdatedEvent := kafka.StockUpdatedEvent{
			BaseEvent: kafka.NewBaseEvent(kafka.EventTypeStockUpdated, "product-service"),
			Data: kafka.StockUpdatedData{
				ProductID: item.ProductID,
				Quantity:  -item.Quantity, // Negative because we're reducing stock
				NewStock:  newStock,
				Reason:    "payment_completed",
				OrderID:   &event.Data.OrderID,
				PaymentID: &event.Data.PaymentID,
			},
		}

		if err := c.eventPublisher.PublishStockUpdated(stockUpdatedEvent); err != nil {
			return fmt.Errorf("failed to publish stock updated event for product %d: %w", item.ProductID, err)
		}

		log.Printf("Stock updated for product %d: -%d units, new stock: %d", 
			item.ProductID, item.Quantity, newStock)
	}

	return nil
}

// HandlePaymentFailed handles payment failed events
func (c *ProductConsumer) HandlePaymentFailed(ctx context.Context, event kafka.PaymentFailedEvent) error {
	log.Printf("Processing payment failed event for stock restoration: %s", event.Data.PaymentID)

	// In case of payment failure, we might need to restore reserved stock
	// This would depend on your business logic - whether you reserve stock before payment
	// or only update stock after successful payment

	log.Printf("Payment failed, no stock changes needed for payment %s", event.Data.PaymentID)
	return nil
}

// HandlePaymentCancelled handles payment cancelled events
func (c *ProductConsumer) HandlePaymentCancelled(ctx context.Context, event kafka.PaymentCancelledEvent) error {
	log.Printf("Processing payment cancelled event for stock restoration: %s", event.Data.PaymentID)

	// Similar to payment failed, restore any reserved stock
	log.Printf("Payment cancelled, no stock changes needed for payment %s", event.Data.PaymentID)
	return nil
}

// HandleStockUpdated handles stock updated events (from other services)
func (c *ProductConsumer) HandleStockUpdated(ctx context.Context, event kafka.StockUpdatedEvent) error {
	log.Printf("Processing stock updated event for product %d: %+d units, new stock: %d", 
		event.Data.ProductID, event.Data.Quantity, event.Data.NewStock)

	// Here you would update the actual stock in your product database
	// This is where the real stock update logic would go

	return nil
}

// Start starts the product consumer
func (c *ProductConsumer) Start(ctx context.Context) error {
	log.Println("Starting product consumer...")

	// In a real implementation, you would:
	// 1. Create a Kafka consumer
	// 2. Subscribe to payment and stock events
	// 3. Handle events in a loop

	// For now, this is a placeholder
	log.Println("Product consumer started successfully")
	return nil
}

// Stop stops the product consumer
func (c *ProductConsumer) Stop() error {
	log.Println("Stopping product consumer...")
	return nil
}
