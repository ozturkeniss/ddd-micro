package monitoring

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics holds all the prometheus metrics for payment service
type PrometheusMetrics struct {
	HTTPRequestsTotal         *prometheus.CounterVec
	HTTPRequestDuration       *prometheus.HistogramVec
	HTTPRequestsInFlight      *prometheus.GaugeVec
	PaymentCreations          prometheus.Counter
	PaymentProcessing         prometheus.Counter
	PaymentCompletions        prometheus.Counter
	PaymentCancellations      prometheus.Counter
	PaymentFailures           prometheus.Counter
	RefundCreations           prometheus.Counter
	RefundProcessing          prometheus.Counter
	RefundCompletions         prometheus.Counter
	RefundFailures            prometheus.Counter
	PaymentMethodAdditions    prometheus.Counter
	PaymentMethodUpdates      prometheus.Counter
	PaymentMethodDeletions    prometheus.Counter
	ActivePayments            prometheus.Gauge
	TotalPaymentAmount        prometheus.Gauge
	AveragePaymentAmount      prometheus.Gauge
	PaymentProcessingDuration *prometheus.HistogramVec
	DatabaseConnections       prometheus.Gauge
	DatabaseQueryDuration     *prometheus.HistogramVec
	ExternalAPICalls          *prometheus.CounterVec
	ExternalAPIDuration       *prometheus.HistogramVec
	StripeAPICalls            *prometheus.CounterVec
	StripeAPIDuration         *prometheus.HistogramVec
	KafkaMessagesPublished    *prometheus.CounterVec
	KafkaMessagesConsumed     *prometheus.CounterVec
	CacheHits                 prometheus.Counter
	CacheMisses               prometheus.Counter
}

// NewPrometheusMetrics creates a new instance of PrometheusMetrics
func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payment_service_http_requests_total",
				Help: "Total number of HTTP requests to payment service",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "payment_service_http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestsInFlight: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "payment_service_http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
			[]string{"method", "endpoint"},
		),
		PaymentCreations: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_payment_creations_total",
				Help: "Total number of payment creations",
			},
		),
		PaymentProcessing: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_payment_processing_total",
				Help: "Total number of payment processing attempts",
			},
		),
		PaymentCompletions: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_payment_completions_total",
				Help: "Total number of successful payment completions",
			},
		),
		PaymentCancellations: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_payment_cancellations_total",
				Help: "Total number of payment cancellations",
			},
		),
		PaymentFailures: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_payment_failures_total",
				Help: "Total number of payment failures",
			},
		),
		RefundCreations: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_refund_creations_total",
				Help: "Total number of refund creations",
			},
		),
		RefundProcessing: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_refund_processing_total",
				Help: "Total number of refund processing attempts",
			},
		),
		RefundCompletions: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_refund_completions_total",
				Help: "Total number of successful refund completions",
			},
		),
		RefundFailures: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_refund_failures_total",
				Help: "Total number of refund failures",
			},
		),
		PaymentMethodAdditions: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_payment_method_additions_total",
				Help: "Total number of payment method additions",
			},
		),
		PaymentMethodUpdates: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_payment_method_updates_total",
				Help: "Total number of payment method updates",
			},
		),
		PaymentMethodDeletions: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_payment_method_deletions_total",
				Help: "Total number of payment method deletions",
			},
		),
		ActivePayments: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "payment_service_active_payments",
				Help: "Current number of active payments",
			},
		),
		TotalPaymentAmount: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "payment_service_total_payment_amount",
				Help: "Total amount of all payments",
			},
		),
		AveragePaymentAmount: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "payment_service_average_payment_amount",
				Help: "Average payment amount",
			},
		),
		PaymentProcessingDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "payment_service_payment_processing_duration_seconds",
				Help:    "Duration of payment processing in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"payment_method", "status"},
		),
		DatabaseConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "payment_service_database_connections_active",
				Help: "Current number of active database connections",
			},
		),
		DatabaseQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "payment_service_database_query_duration_seconds",
				Help:    "Duration of database queries in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		),
		ExternalAPICalls: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payment_service_external_api_calls_total",
				Help: "Total number of external API calls",
			},
			[]string{"service", "endpoint", "status"},
		),
		ExternalAPIDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "payment_service_external_api_duration_seconds",
				Help:    "Duration of external API calls in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "endpoint"},
		),
		StripeAPICalls: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payment_service_stripe_api_calls_total",
				Help: "Total number of Stripe API calls",
			},
			[]string{"operation", "status"},
		),
		StripeAPIDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "payment_service_stripe_api_duration_seconds",
				Help:    "Duration of Stripe API calls in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation"},
		),
		KafkaMessagesPublished: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payment_service_kafka_messages_published_total",
				Help: "Total number of Kafka messages published",
			},
			[]string{"topic", "status"},
		),
		KafkaMessagesConsumed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payment_service_kafka_messages_consumed_total",
				Help: "Total number of Kafka messages consumed",
			},
			[]string{"topic", "status"},
		),
		CacheHits: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_cache_hits_total",
				Help: "Total number of cache hits",
			},
		),
		CacheMisses: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_service_cache_misses_total",
				Help: "Total number of cache misses",
			},
		),
	}
}

// PrometheusMiddleware returns a Gin middleware for Prometheus metrics
func PrometheusMiddleware(metrics *PrometheusMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Increment requests in flight
		metrics.HTTPRequestsInFlight.WithLabelValues(c.Request.Method, c.FullPath()).Inc()
		defer metrics.HTTPRequestsInFlight.WithLabelValues(c.Request.Method, c.FullPath()).Dec()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		
		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			string(rune(status)),
		).Inc()
		
		metrics.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			string(rune(status)),
		).Observe(duration)
	}
}

// RecordPaymentCreation increments the payment creation counter
func (m *PrometheusMetrics) RecordPaymentCreation() {
	m.PaymentCreations.Inc()
}

// RecordPaymentProcessing increments the payment processing counter
func (m *PrometheusMetrics) RecordPaymentProcessing() {
	m.PaymentProcessing.Inc()
}

// RecordPaymentCompletion increments the payment completion counter
func (m *PrometheusMetrics) RecordPaymentCompletion() {
	m.PaymentCompletions.Inc()
}

// RecordPaymentCancellation increments the payment cancellation counter
func (m *PrometheusMetrics) RecordPaymentCancellation() {
	m.PaymentCancellations.Inc()
}

// RecordPaymentFailure increments the payment failure counter
func (m *PrometheusMetrics) RecordPaymentFailure() {
	m.PaymentFailures.Inc()
}

// RecordRefundCreation increments the refund creation counter
func (m *PrometheusMetrics) RecordRefundCreation() {
	m.RefundCreations.Inc()
}

// RecordRefundProcessing increments the refund processing counter
func (m *PrometheusMetrics) RecordRefundProcessing() {
	m.RefundProcessing.Inc()
}

// RecordRefundCompletion increments the refund completion counter
func (m *PrometheusMetrics) RecordRefundCompletion() {
	m.RefundCompletions.Inc()
}

// RecordRefundFailure increments the refund failure counter
func (m *PrometheusMetrics) RecordRefundFailure() {
	m.RefundFailures.Inc()
}

// RecordPaymentMethodAddition increments the payment method addition counter
func (m *PrometheusMetrics) RecordPaymentMethodAddition() {
	m.PaymentMethodAdditions.Inc()
}

// RecordPaymentMethodUpdate increments the payment method update counter
func (m *PrometheusMetrics) RecordPaymentMethodUpdate() {
	m.PaymentMethodUpdates.Inc()
}

// RecordPaymentMethodDeletion increments the payment method deletion counter
func (m *PrometheusMetrics) RecordPaymentMethodDeletion() {
	m.PaymentMethodDeletions.Inc()
}

// SetActivePayments sets the active payments gauge
func (m *PrometheusMetrics) SetActivePayments(count float64) {
	m.ActivePayments.Set(count)
}

// SetTotalPaymentAmount sets the total payment amount gauge
func (m *PrometheusMetrics) SetTotalPaymentAmount(amount float64) {
	m.TotalPaymentAmount.Set(amount)
}

// SetAveragePaymentAmount sets the average payment amount gauge
func (m *PrometheusMetrics) SetAveragePaymentAmount(amount float64) {
	m.AveragePaymentAmount.Set(amount)
}

// RecordPaymentProcessingDuration records the duration of payment processing
func (m *PrometheusMetrics) RecordPaymentProcessingDuration(paymentMethod, status string, duration time.Duration) {
	m.PaymentProcessingDuration.WithLabelValues(paymentMethod, status).Observe(duration.Seconds())
}

// SetDatabaseConnections sets the database connections gauge
func (m *PrometheusMetrics) SetDatabaseConnections(count float64) {
	m.DatabaseConnections.Set(count)
}

// RecordDatabaseQueryDuration records the duration of a database query
func (m *PrometheusMetrics) RecordDatabaseQueryDuration(operation, table string, duration time.Duration) {
	m.DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// RecordExternalAPICall records an external API call
func (m *PrometheusMetrics) RecordExternalAPICall(service, endpoint, status string) {
	m.ExternalAPICalls.WithLabelValues(service, endpoint, status).Inc()
}

// RecordExternalAPIDuration records the duration of an external API call
func (m *PrometheusMetrics) RecordExternalAPIDuration(service, endpoint string, duration time.Duration) {
	m.ExternalAPIDuration.WithLabelValues(service, endpoint).Observe(duration.Seconds())
}

// RecordStripeAPICall records a Stripe API call
func (m *PrometheusMetrics) RecordStripeAPICall(operation, status string) {
	m.StripeAPICalls.WithLabelValues(operation, status).Inc()
}

// RecordStripeAPIDuration records the duration of a Stripe API call
func (m *PrometheusMetrics) RecordStripeAPIDuration(operation string, duration time.Duration) {
	m.StripeAPIDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordKafkaMessagePublished records a published Kafka message
func (m *PrometheusMetrics) RecordKafkaMessagePublished(topic, status string) {
	m.KafkaMessagesPublished.WithLabelValues(topic, status).Inc()
}

// RecordKafkaMessageConsumed records a consumed Kafka message
func (m *PrometheusMetrics) RecordKafkaMessageConsumed(topic, status string) {
	m.KafkaMessagesConsumed.WithLabelValues(topic, status).Inc()
}

// RecordCacheHit increments the cache hit counter
func (m *PrometheusMetrics) RecordCacheHit() {
	m.CacheHits.Inc()
}

// RecordCacheMiss increments the cache miss counter
func (m *PrometheusMetrics) RecordCacheMiss() {
	m.CacheMisses.Inc()
}
