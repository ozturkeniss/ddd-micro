package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// kafkaConsumer implements EventConsumer interface
type kafkaConsumer struct {
	consumer sarama.ConsumerGroup
	config   *ConsumerConfig
	handlers map[EventType]func([]byte) error
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// ConsumerConfig holds configuration for the Kafka consumer
type ConsumerConfig struct {
	Brokers       []string
	Topic         string
	GroupID       string
	Offset        int64
	RetryAttempts int
	RetryDelay    time.Duration
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(config *ConsumerConfig) (EventConsumer, error) {
	// Configure consumer
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = config.Offset
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Group.Session.Timeout = 10 * time.Second
	saramaConfig.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	// Create consumer group
	consumer, err := sarama.NewConsumerGroup(config.Brokers, config.GroupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &kafkaConsumer{
		consumer: consumer,
		config:   config,
		handlers: make(map[EventType]func([]byte) error),
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// ConsumePaymentCompleted registers a handler for payment completed events
func (c *kafkaConsumer) ConsumePaymentCompleted(handler func(PaymentCompletedEvent) error) error {
	c.handlers[EventTypePaymentCompleted] = func(data []byte) error {
		var event PaymentCompletedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to unmarshal payment completed event: %w", err)
		}
		return handler(event)
	}
	return nil
}

// ConsumePaymentFailed registers a handler for payment failed events
func (c *kafkaConsumer) ConsumePaymentFailed(handler func(PaymentFailedEvent) error) error {
	c.handlers[EventTypePaymentFailed] = func(data []byte) error {
		var event PaymentFailedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to unmarshal payment failed event: %w", err)
		}
		return handler(event)
	}
	return nil
}

// ConsumePaymentCancelled registers a handler for payment cancelled events
func (c *kafkaConsumer) ConsumePaymentCancelled(handler func(PaymentCancelledEvent) error) error {
	c.handlers[EventTypePaymentCancelled] = func(data []byte) error {
		var event PaymentCancelledEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to unmarshal payment cancelled event: %w", err)
		}
		return handler(event)
	}
	return nil
}

// ConsumeStockUpdated registers a handler for stock updated events
func (c *kafkaConsumer) ConsumeStockUpdated(handler func(StockUpdatedEvent) error) error {
	c.handlers[EventTypeStockUpdated] = func(data []byte) error {
		var event StockUpdatedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to unmarshal stock updated event: %w", err)
		}
		return handler(event)
	}
	return nil
}

// ConsumeBasketCleared registers a handler for basket cleared events
func (c *kafkaConsumer) ConsumeBasketCleared(handler func(BasketClearedEvent) error) error {
	c.handlers[EventTypeBasketCleared] = func(data []byte) error {
		var event BasketClearedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to unmarshal basket cleared event: %w", err)
		}
		return handler(event)
	}
	return nil
}

// ConsumeOrderCreated registers a handler for order created events
func (c *kafkaConsumer) ConsumeOrderCreated(handler func(OrderCreatedEvent) error) error {
	c.handlers[EventTypeOrderCreated] = func(data []byte) error {
		var event OrderCreatedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to unmarshal order created event: %w", err)
		}
		return handler(event)
	}
	return nil
}

// Start starts the consumer
func (c *kafkaConsumer) Start() error {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				log.Println("Consumer context cancelled")
				return
			default:
				// Start consuming
				if err := c.consumer.Consume(c.ctx, []string{c.config.Topic}, c); err != nil {
					log.Printf("Error from consumer: %v", err)
					time.Sleep(c.config.RetryDelay)
				}
			}
		}
	}()

	// Start error handler
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				return
			case err := <-c.consumer.Errors():
				log.Printf("Consumer error: %v", err)
			}
		}
	}()

	log.Printf("Kafka consumer started for topic: %s, group: %s", c.config.Topic, c.config.GroupID)
	return nil
}

// Stop stops the consumer
func (c *kafkaConsumer) Stop() error {
	c.cancel()
	c.wg.Wait()
	return c.consumer.Close()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *kafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session setup")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *kafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (c *kafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			// Get event type from headers
			eventType := EventType("")
			for _, header := range message.Headers {
				if string(header.Key) == "event-type" {
					eventType = EventType(header.Value)
					break
				}
			}

			// Process message
			if handler, exists := c.handlers[eventType]; exists {
				if err := c.processMessage(handler, message.Value); err != nil {
					log.Printf("Error processing message: %v", err)
					// In production, you might want to implement retry logic or dead letter queue
				}
			} else {
				log.Printf("No handler found for event type: %s", eventType)
			}

			// Mark message as processed
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage processes a single message with retry logic
func (c *kafkaConsumer) processMessage(handler func([]byte) error, data []byte) error {
	var lastErr error
	for attempt := 0; attempt < c.config.RetryAttempts; attempt++ {
		if err := handler(data); err != nil {
			lastErr = err
			if attempt < c.config.RetryAttempts-1 {
				log.Printf("Handler failed (attempt %d/%d): %v, retrying in %v",
					attempt+1, c.config.RetryAttempts, err, c.config.RetryDelay)
				time.Sleep(c.config.RetryDelay)
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("handler failed after %d attempts: %w", c.config.RetryAttempts, lastErr)
}

// ConsumerGroup represents a consumer group
type ConsumerGroup struct {
	consumer EventConsumer
	handlers map[EventType]func([]byte) error
}

// NewConsumerGroup creates a new consumer group
func NewConsumerGroup(consumer EventConsumer) *ConsumerGroup {
	return &ConsumerGroup{
		consumer: consumer,
		handlers: make(map[EventType]func([]byte) error),
	}
}

// RegisterHandler registers an event handler
func (cg *ConsumerGroup) RegisterHandler(eventType EventType, handler func([]byte) error) {
	cg.handlers[eventType] = handler
}

// Start starts the consumer group
func (cg *ConsumerGroup) Start() error {
	return cg.consumer.Start()
}

// Stop stops the consumer group
func (cg *ConsumerGroup) Stop() error {
	return cg.consumer.Stop()
}
