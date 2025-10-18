package kafka

import (
	"encoding/json"
	"time"
)

// EventType represents the type of event
type EventType string

const (
	EventTypePaymentCompleted EventType = "payment.completed"
	EventTypePaymentFailed    EventType = "payment.failed"
	EventTypePaymentCancelled EventType = "payment.cancelled"
	EventTypeStockUpdated     EventType = "stock.updated"
	EventTypeBasketCleared    EventType = "basket.cleared"
	EventTypeOrderCreated     EventType = "order.created"
)

// BaseEvent represents the base structure for all events
type BaseEvent struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
	Version   string    `json:"version"`
}

// PaymentCompletedEvent represents a payment completion event
type PaymentCompletedEvent struct {
	BaseEvent
	Data PaymentCompletedData `json:"data"`
}

// PaymentCompletedData contains the payment completion data
type PaymentCompletedData struct {
	PaymentID     string                 `json:"payment_id"`
	UserID        uint                   `json:"user_id"`
	OrderID       string                 `json:"order_id"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	PaymentMethod string                 `json:"payment_method"`
	Items         []PaymentItem          `json:"items"`
	BasketID      *string                `json:"basket_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// PaymentItem represents an item in the payment
type PaymentItem struct {
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

// PaymentFailedEvent represents a payment failure event
type PaymentFailedEvent struct {
	BaseEvent
	Data PaymentFailedData `json:"data"`
}

// PaymentFailedData contains the payment failure data
type PaymentFailedData struct {
	PaymentID     string `json:"payment_id"`
	UserID        uint   `json:"user_id"`
	OrderID       string `json:"order_id"`
	Amount        float64 `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	Reason        string `json:"reason"`
	BasketID      *string `json:"basket_id,omitempty"`
}

// PaymentCancelledEvent represents a payment cancellation event
type PaymentCancelledEvent struct {
	BaseEvent
	Data PaymentCancelledData `json:"data"`
}

// PaymentCancelledData contains the payment cancellation data
type PaymentCancelledData struct {
	PaymentID     string `json:"payment_id"`
	UserID        uint   `json:"user_id"`
	OrderID       string `json:"order_id"`
	Amount        float64 `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	Reason        string `json:"reason"`
	BasketID      *string `json:"basket_id,omitempty"`
}

// StockUpdatedEvent represents a stock update event
type StockUpdatedEvent struct {
	BaseEvent
	Data StockUpdatedData `json:"data"`
}

// StockUpdatedData contains the stock update data
type StockUpdatedData struct {
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"` // Positive for increase, negative for decrease
	NewStock  int     `json:"new_stock"`
	Reason    string  `json:"reason"`
	OrderID   *string `json:"order_id,omitempty"`
	PaymentID *string `json:"payment_id,omitempty"`
}

// BasketClearedEvent represents a basket clearing event
type BasketClearedEvent struct {
	BaseEvent
	Data BasketClearedData `json:"data"`
}

// BasketClearedData contains the basket clearing data
type BasketClearedData struct {
	UserID   uint           `json:"user_id"`
	BasketID string         `json:"basket_id"`
	Items    []PaymentItem  `json:"items"`
	Reason   string         `json:"reason"`
	OrderID  *string        `json:"order_id,omitempty"`
	PaymentID *string       `json:"payment_id,omitempty"`
}

// OrderCreatedEvent represents an order creation event
type OrderCreatedEvent struct {
	BaseEvent
	Data OrderCreatedData `json:"data"`
}

// OrderCreatedData contains the order creation data
type OrderCreatedData struct {
	OrderID       string        `json:"order_id"`
	UserID        uint          `json:"user_id"`
	PaymentID     string        `json:"payment_id"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	Items         []PaymentItem `json:"items"`
	ShippingInfo  ShippingInfo  `json:"shipping_info"`
	BillingInfo   BillingInfo   `json:"billing_info"`
}

// ShippingInfo represents shipping information
type ShippingInfo struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
	Phone   string `json:"phone"`
}

// BillingInfo represents billing information
type BillingInfo struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

// EventHandler defines the interface for event handlers
type EventHandler interface {
	Handle(event BaseEvent) error
}

// EventPublisher defines the interface for event publishers
type EventPublisher interface {
	PublishPaymentCompleted(event PaymentCompletedEvent) error
	PublishPaymentFailed(event PaymentFailedEvent) error
	PublishPaymentCancelled(event PaymentCancelledEvent) error
	PublishStockUpdated(event StockUpdatedEvent) error
	PublishBasketCleared(event BasketClearedEvent) error
	PublishOrderCreated(event OrderCreatedEvent) error
}

// EventConsumer defines the interface for event consumers
type EventConsumer interface {
	ConsumePaymentCompleted(handler func(PaymentCompletedEvent) error) error
	ConsumePaymentFailed(handler func(PaymentFailedEvent) error) error
	ConsumePaymentCancelled(handler func(PaymentCancelledEvent) error) error
	ConsumeStockUpdated(handler func(StockUpdatedEvent) error) error
	ConsumeBasketCleared(handler func(BasketClearedEvent) error) error
	ConsumeOrderCreated(handler func(OrderCreatedEvent) error) error
	Start() error
	Stop() error
}

// Helper functions

// NewBaseEvent creates a new base event
func NewBaseEvent(eventType EventType, source string) BaseEvent {
	return BaseEvent{
		ID:        generateEventID(),
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Source:    source,
		Version:   "1.0",
	}
}

// ToJSON converts an event to JSON
func (e BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON converts JSON to an event
func FromJSON(data []byte, event interface{}) error {
	return json.Unmarshal(data, event)
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
