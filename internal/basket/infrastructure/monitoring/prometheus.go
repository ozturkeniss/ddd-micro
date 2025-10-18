package monitoring

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics holds all the prometheus metrics for basket service
type PrometheusMetrics struct {
	HTTPRequestsTotal       *prometheus.CounterVec
	HTTPRequestDuration     *prometheus.HistogramVec
	HTTPRequestsInFlight    *prometheus.GaugeVec
	BasketCreations         prometheus.Counter
	BasketRetrievals        prometheus.Counter
	BasketUpdates           prometheus.Counter
	BasketDeletions         prometheus.Counter
	BasketClearings         prometheus.Counter
	ItemAdditions           prometheus.Counter
	ItemUpdates             prometheus.Counter
	ItemRemovals            prometheus.Counter
	BasketViews             prometheus.Counter
	ActiveBaskets           prometheus.Gauge
	TotalItemsInBaskets     prometheus.Gauge
	BasketExpirations       prometheus.Counter
	BasketCleanups          prometheus.Counter
	RedisOperations         *prometheus.CounterVec
	RedisOperationDuration  *prometheus.HistogramVec
	ExternalAPICalls        *prometheus.CounterVec
	ExternalAPIDuration     *prometheus.HistogramVec
	CacheHits               prometheus.Counter
	CacheMisses             prometheus.Counter
}

// NewPrometheusMetrics creates a new instance of PrometheusMetrics
func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "basket_service_http_requests_total",
				Help: "Total number of HTTP requests to basket service",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "basket_service_http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestsInFlight: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "basket_service_http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
			[]string{"method", "endpoint"},
		),
		BasketCreations: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_basket_creations_total",
				Help: "Total number of basket creations",
			},
		),
		BasketRetrievals: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_basket_retrievals_total",
				Help: "Total number of basket retrievals",
			},
		),
		BasketUpdates: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_basket_updates_total",
				Help: "Total number of basket updates",
			},
		),
		BasketDeletions: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_basket_deletions_total",
				Help: "Total number of basket deletions",
			},
		),
		BasketClearings: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_basket_clearings_total",
				Help: "Total number of basket clearings",
			},
		),
		ItemAdditions: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_item_additions_total",
				Help: "Total number of item additions to baskets",
			},
		),
		ItemUpdates: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_item_updates_total",
				Help: "Total number of item updates in baskets",
			},
		),
		ItemRemovals: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_item_removals_total",
				Help: "Total number of item removals from baskets",
			},
		),
		BasketViews: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_basket_views_total",
				Help: "Total number of basket views",
			},
		),
		ActiveBaskets: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "basket_service_active_baskets",
				Help: "Current number of active baskets",
			},
		),
		TotalItemsInBaskets: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "basket_service_total_items_in_baskets",
				Help: "Current total number of items in all baskets",
			},
		),
		BasketExpirations: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_basket_expirations_total",
				Help: "Total number of basket expirations",
			},
		),
		BasketCleanups: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_basket_cleanups_total",
				Help: "Total number of basket cleanups",
			},
		),
		RedisOperations: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "basket_service_redis_operations_total",
				Help: "Total number of Redis operations",
			},
			[]string{"operation", "status"},
		),
		RedisOperationDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "basket_service_redis_operation_duration_seconds",
				Help:    "Duration of Redis operations in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation"},
		),
		ExternalAPICalls: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "basket_service_external_api_calls_total",
				Help: "Total number of external API calls",
			},
			[]string{"service", "endpoint", "status"},
		),
		ExternalAPIDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "basket_service_external_api_duration_seconds",
				Help:    "Duration of external API calls in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "endpoint"},
		),
		CacheHits: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_cache_hits_total",
				Help: "Total number of cache hits",
			},
		),
		CacheMisses: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "basket_service_cache_misses_total",
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

// RecordBasketCreation increments the basket creation counter
func (m *PrometheusMetrics) RecordBasketCreation() {
	m.BasketCreations.Inc()
}

// RecordBasketRetrieval increments the basket retrieval counter
func (m *PrometheusMetrics) RecordBasketRetrieval() {
	m.BasketRetrievals.Inc()
}

// RecordBasketUpdate increments the basket update counter
func (m *PrometheusMetrics) RecordBasketUpdate() {
	m.BasketUpdates.Inc()
}

// RecordBasketDeletion increments the basket deletion counter
func (m *PrometheusMetrics) RecordBasketDeletion() {
	m.BasketDeletions.Inc()
}

// RecordBasketClearing increments the basket clearing counter
func (m *PrometheusMetrics) RecordBasketClearing() {
	m.BasketClearings.Inc()
}

// RecordItemAddition increments the item addition counter
func (m *PrometheusMetrics) RecordItemAddition() {
	m.ItemAdditions.Inc()
}

// RecordItemUpdate increments the item update counter
func (m *PrometheusMetrics) RecordItemUpdate() {
	m.ItemUpdates.Inc()
}

// RecordItemRemoval increments the item removal counter
func (m *PrometheusMetrics) RecordItemRemoval() {
	m.ItemRemovals.Inc()
}

// RecordBasketView increments the basket view counter
func (m *PrometheusMetrics) RecordBasketView() {
	m.BasketViews.Inc()
}

// SetActiveBaskets sets the active baskets gauge
func (m *PrometheusMetrics) SetActiveBaskets(count float64) {
	m.ActiveBaskets.Set(count)
}

// SetTotalItemsInBaskets sets the total items in baskets gauge
func (m *PrometheusMetrics) SetTotalItemsInBaskets(count float64) {
	m.TotalItemsInBaskets.Set(count)
}

// RecordBasketExpiration increments the basket expiration counter
func (m *PrometheusMetrics) RecordBasketExpiration() {
	m.BasketExpirations.Inc()
}

// RecordBasketCleanup increments the basket cleanup counter
func (m *PrometheusMetrics) RecordBasketCleanup() {
	m.BasketCleanups.Inc()
}

// RecordRedisOperation records a Redis operation
func (m *PrometheusMetrics) RecordRedisOperation(operation, status string) {
	m.RedisOperations.WithLabelValues(operation, status).Inc()
}

// RecordRedisOperationDuration records the duration of a Redis operation
func (m *PrometheusMetrics) RecordRedisOperationDuration(operation string, duration time.Duration) {
	m.RedisOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordExternalAPICall records an external API call
func (m *PrometheusMetrics) RecordExternalAPICall(service, endpoint, status string) {
	m.ExternalAPICalls.WithLabelValues(service, endpoint, status).Inc()
}

// RecordExternalAPIDuration records the duration of an external API call
func (m *PrometheusMetrics) RecordExternalAPIDuration(service, endpoint string, duration time.Duration) {
	m.ExternalAPIDuration.WithLabelValues(service, endpoint).Observe(duration.Seconds())
}

// RecordCacheHit increments the cache hit counter
func (m *PrometheusMetrics) RecordCacheHit() {
	m.CacheHits.Inc()
}

// RecordCacheMiss increments the cache miss counter
func (m *PrometheusMetrics) RecordCacheMiss() {
	m.CacheMisses.Inc()
}
