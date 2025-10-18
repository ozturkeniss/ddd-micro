package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// kafkaPublisher implements EventPublisher interface
type kafkaPublisher struct {
	producer sarama.SyncProducer
	config   *PublisherConfig
}

// PublisherConfig holds configuration for the Kafka publisher
type PublisherConfig struct {
	Brokers []string
	Topic   string
	Timeout time.Duration
}

// NewKafkaPublisher creates a new Kafka publisher
func NewKafkaPublisher(config *PublisherConfig) (EventPublisher, error) {
	// Configure producer
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Timeout = config.Timeout

	// Create producer
	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &kafkaPublisher{
		producer: producer,
		config:   config,
	}, nil
}

// PublishPaymentCompleted publishes a payment completed event
func (p *kafkaPublisher) PublishPaymentCompleted(event PaymentCompletedEvent) error {
	return p.publishEvent(event.BaseEvent.Type, event)
}

// PublishPaymentFailed publishes a payment failed event
func (p *kafkaPublisher) PublishPaymentFailed(event PaymentFailedEvent) error {
	return p.publishEvent(event.BaseEvent.Type, event)
}

// PublishPaymentCancelled publishes a payment cancelled event
func (p *kafkaPublisher) PublishPaymentCancelled(event PaymentCancelledEvent) error {
	return p.publishEvent(event.BaseEvent.Type, event)
}

// PublishStockUpdated publishes a stock updated event
func (p *kafkaPublisher) PublishStockUpdated(event StockUpdatedEvent) error {
	return p.publishEvent(event.BaseEvent.Type, event)
}

// PublishBasketCleared publishes a basket cleared event
func (p *kafkaPublisher) PublishBasketCleared(event BasketClearedEvent) error {
	return p.publishEvent(event.BaseEvent.Type, event)
}

// PublishOrderCreated publishes an order created event
func (p *kafkaPublisher) PublishOrderCreated(event OrderCreatedEvent) error {
	return p.publishEvent(event.BaseEvent.Type, event)
}

// publishEvent publishes a generic event to Kafka
func (p *kafkaPublisher) publishEvent(eventType EventType, event interface{}) error {
	// Serialize event to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create message
	message := &sarama.ProducerMessage{
		Topic: p.config.Topic,
		Key:   sarama.StringEncoder(eventType),
		Value: sarama.ByteEncoder(eventData),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event-type"),
				Value: []byte(eventType),
			},
			{
				Key:   []byte("timestamp"),
				Value: []byte(time.Now().UTC().Format(time.RFC3339)),
			},
		},
	}

	// Send message
	partition, offset, err := p.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	log.Printf("Event published successfully: type=%s, partition=%d, offset=%d", eventType, partition, offset)
	return nil
}

// Close closes the publisher
func (p *kafkaPublisher) Close() error {
	return p.producer.Close()
}

// PaymentEventPublisher is a specialized publisher for payment events
type PaymentEventPublisher struct {
	publisher EventPublisher
}

// NewPaymentEventPublisher creates a new payment event publisher
func NewPaymentEventPublisher(publisher EventPublisher) *PaymentEventPublisher {
	return &PaymentEventPublisher{
		publisher: publisher,
	}
}

// PublishPaymentCompleted publishes a payment completed event with basket clearing
func (p *PaymentEventPublisher) PublishPaymentCompleted(ctx context.Context, paymentID string, userID uint, orderID string, amount float64, currency string, paymentMethod string, items []PaymentItem, basketID *string) error {
	event := PaymentCompletedEvent{
		BaseEvent: NewBaseEvent(EventTypePaymentCompleted, "payment-service"),
		Data: PaymentCompletedData{
			PaymentID:     paymentID,
			UserID:        userID,
			OrderID:       orderID,
			Amount:        amount,
			Currency:      currency,
			PaymentMethod: paymentMethod,
			Items:         items,
			BasketID:      basketID,
			Metadata: map[string]interface{}{
				"timestamp": time.Now().UTC(),
				"source":    "payment-service",
			},
		},
	}

	return p.publisher.PublishPaymentCompleted(event)
}

// PublishPaymentFailed publishes a payment failed event
func (p *PaymentEventPublisher) PublishPaymentFailed(ctx context.Context, paymentID string, userID uint, orderID string, amount float64, currency string, paymentMethod string, reason string, basketID *string) error {
	event := PaymentFailedEvent{
		BaseEvent: NewBaseEvent(EventTypePaymentFailed, "payment-service"),
		Data: PaymentFailedData{
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
	event := PaymentCancelledEvent{
		BaseEvent: NewBaseEvent(EventTypePaymentCancelled, "payment-service"),
		Data: PaymentCancelledData{
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
	event := StockUpdatedEvent{
		BaseEvent: NewBaseEvent(EventTypeStockUpdated, "payment-service"),
		Data: StockUpdatedData{
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
func (p *PaymentEventPublisher) PublishBasketCleared(ctx context.Context, userID uint, basketID string, items []PaymentItem, reason string, orderID *string, paymentID *string) error {
	event := BasketClearedEvent{
		BaseEvent: NewBaseEvent(EventTypeBasketCleared, "payment-service"),
		Data: BasketClearedData{
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
