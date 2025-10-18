package kafka

import (
	"context"
	"fmt"

	"github.com/ddd-micro/kafka"
)

// PaymentEventPublisher handles payment-related Kafka events
type PaymentEventPublisher struct {
	publisher kafka.EventPublisher
}

// NewPaymentEventPublisher creates a new payment event publisher
func NewPaymentEventPublisher(publisher kafka.EventPublisher) *PaymentEventPublisher {
	return &PaymentEventPublisher{
		publisher: publisher,
	}
}

// PublishPaymentCompleted publishes a payment completed event
func (p *PaymentEventPublisher) PublishPaymentCompleted(ctx context.Context, paymentID string, userID uint, orderID string, amount float64, currency string, paymentMethod string, items []kafka.PaymentItem, basketID *string) error {
	event := kafka.PaymentCompletedEvent{
		BaseEvent: kafka.NewBaseEvent(kafka.EventTypePaymentCompleted, "payment-service"),
		Data: kafka.PaymentCompletedData{
			PaymentID:     paymentID,
			UserID:        userID,
			OrderID:       orderID,
			Amount:        amount,
			Currency:      currency,
			PaymentMethod: paymentMethod,
			Items:         items,
			BasketID:      basketID,
			Metadata: map[string]interface{}{
				"timestamp": "2024-01-01T00:00:00Z", // This should be actual timestamp
				"source":    "payment-service",
			},
		},
	}

	return p.publisher.PublishPaymentCompleted(event)
}

// PublishPaymentFailed publishes a payment failed event
func (p *PaymentEventPublisher) PublishPaymentFailed(ctx context.Context, paymentID string, userID uint, orderID string, amount float64, currency string, paymentMethod string, reason string, basketID *string) error {
	event := kafka.PaymentFailedEvent{
		BaseEvent: kafka.NewBaseEvent(kafka.EventTypePaymentFailed, "payment-service"),
		Data: kafka.PaymentFailedData{
			PaymentID:     paymentID,
			UserID:        userID,
			OrderID:       orderID,
			Amount:        amount,
			Currency:      currency,
			PaymentMethod: paymentMethod,
			Reason:        reason,
			BasketID:      basketID,
		},
	}

	return p.publisher.PublishPaymentFailed(event)
}

// PublishPaymentCancelled publishes a payment cancelled event
func (p *PaymentEventPublisher) PublishPaymentCancelled(ctx context.Context, paymentID string, userID uint, orderID string, amount float64, currency string, paymentMethod string, reason string, basketID *string) error {
	event := kafka.PaymentCancelledEvent{
		BaseEvent: kafka.NewBaseEvent(kafka.EventTypePaymentCancelled, "payment-service"),
		Data: kafka.PaymentCancelledData{
			PaymentID:     paymentID,
			UserID:        userID,
			OrderID:       orderID,
			Amount:        amount,
			Currency:      currency,
			PaymentMethod: paymentMethod,
			Reason:        reason,
			BasketID:      basketID,
		},
	}

	return p.publisher.PublishPaymentCancelled(event)
}

// PublishStockUpdated publishes a stock updated event
func (p *PaymentEventPublisher) PublishStockUpdated(ctx context.Context, productID uint, quantity int, newStock int, reason string, orderID *string, paymentID *string) error {
	event := kafka.StockUpdatedEvent{
		BaseEvent: kafka.NewBaseEvent(kafka.EventTypeStockUpdated, "payment-service"),
		Data: kafka.StockUpdatedData{
			ProductID: productID,
			Quantity:  quantity,
			NewStock:  newStock,
			Reason:    reason,
			OrderID:   orderID,
			PaymentID: paymentID,
		},
	}

	return p.publisher.PublishStockUpdated(event)
}

// PublishBasketCleared publishes a basket cleared event
func (p *PaymentEventPublisher) PublishBasketCleared(ctx context.Context, userID uint, basketID string, items []kafka.PaymentItem, reason string, orderID *string, paymentID *string) error {
	event := kafka.BasketClearedEvent{
		BaseEvent: kafka.NewBaseEvent(kafka.EventTypeBasketCleared, "payment-service"),
		Data: kafka.BasketClearedData{
			UserID:    userID,
			BasketID:  basketID,
			Items:     items,
			Reason:    reason,
			OrderID:   orderID,
			PaymentID: paymentID,
		},
	}

	return p.publisher.PublishBasketCleared(event)
}
