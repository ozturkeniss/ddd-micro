package monitoring

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics holds all the prometheus metrics for product service
type PrometheusMetrics struct {
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPRequestsInFlight  *prometheus.GaugeVec
	ProductViews          prometheus.Counter
	ProductCreations      prometheus.Counter
	ProductUpdates        prometheus.Counter
	ProductDeletions      prometheus.Counter
	StockUpdates          prometheus.Counter
	StockReductions       prometheus.Counter
	StockIncreases        prometheus.Counter
	ProductSearches       prometheus.Counter
	DatabaseConnections   prometheus.Gauge
	DatabaseQueryDuration *prometheus.HistogramVec
	CacheHits             prometheus.Counter
	CacheMisses           prometheus.Counter
	ExternalAPICalls      *prometheus.CounterVec
	ExternalAPIDuration   *prometheus.HistogramVec
}

// NewPrometheusMetrics creates a new instance of PrometheusMetrics
func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "product_service_http_requests_total",
				Help: "Total number of HTTP requests to product service",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "product_service_http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestsInFlight: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "product_service_http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
			[]string{"method", "endpoint"},
		),
		ProductViews: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_service_product_views_total",
				Help: "Total number of product views",
			},
		),
		ProductCreations: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_service_product_creations_total",
				Help: "Total number of product creations",
			},
		),
		ProductUpdates: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_service_product_updates_total",
				Help: "Total number of product updates",
			},
		),
		ProductDeletions: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_service_product_deletions_total",
				Help: "Total number of product deletions",
			},
		),
		StockUpdates: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_service_stock_updates_total",
				Help: "Total number of stock updates",
			},
		),
		StockReductions: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_service_stock_reductions_total",
				Help: "Total number of stock reductions",
			},
		),
		StockIncreases: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_service_stock_increases_total",
				Help: "Total number of stock increases",
			},
		),
		ProductSearches: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_service_product_searches_total",
				Help: "Total number of product searches",
			},
		),
		DatabaseConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "product_service_database_connections_active",
				Help: "Current number of active database connections",
			},
		),
		DatabaseQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "product_service_database_query_duration_seconds",
				Help:    "Duration of database queries in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		),
		CacheHits: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_service_cache_hits_total",
				Help: "Total number of cache hits",
			},
		),
		CacheMisses: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_service_cache_misses_total",
				Help: "Total number of cache misses",
			},
		),
		ExternalAPICalls: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "product_service_external_api_calls_total",
				Help: "Total number of external API calls",
			},
			[]string{"service", "endpoint", "status"},
		),
		ExternalAPIDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "product_service_external_api_duration_seconds",
				Help:    "Duration of external API calls in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "endpoint"},
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

// RecordProductView increments the product view counter
func (m *PrometheusMetrics) RecordProductView() {
	m.ProductViews.Inc()
}

// RecordProductCreation increments the product creation counter
func (m *PrometheusMetrics) RecordProductCreation() {
	m.ProductCreations.Inc()
}

// RecordProductUpdate increments the product update counter
func (m *PrometheusMetrics) RecordProductUpdate() {
	m.ProductUpdates.Inc()
}

// RecordProductDeletion increments the product deletion counter
func (m *PrometheusMetrics) RecordProductDeletion() {
	m.ProductDeletions.Inc()
}

// RecordStockUpdate increments the stock update counter
func (m *PrometheusMetrics) RecordStockUpdate() {
	m.StockUpdates.Inc()
}

// RecordStockReduction increments the stock reduction counter
func (m *PrometheusMetrics) RecordStockReduction() {
	m.StockReductions.Inc()
}

// RecordStockIncrease increments the stock increase counter
func (m *PrometheusMetrics) RecordStockIncrease() {
	m.StockIncreases.Inc()
}

// RecordProductSearch increments the product search counter
func (m *PrometheusMetrics) RecordProductSearch() {
	m.ProductSearches.Inc()
}

// SetDatabaseConnections sets the database connections gauge
func (m *PrometheusMetrics) SetDatabaseConnections(count float64) {
	m.DatabaseConnections.Set(count)
}

// RecordDatabaseQuery records the duration of a database query
func (m *PrometheusMetrics) RecordDatabaseQuery(operation, table string, duration time.Duration) {
	m.DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// RecordCacheHit increments the cache hit counter
func (m *PrometheusMetrics) RecordCacheHit() {
	m.CacheHits.Inc()
}

// RecordCacheMiss increments the cache miss counter
func (m *PrometheusMetrics) RecordCacheMiss() {
	m.CacheMisses.Inc()
}

// RecordExternalAPICall records an external API call
func (m *PrometheusMetrics) RecordExternalAPICall(service, endpoint, status string) {
	m.ExternalAPICalls.WithLabelValues(service, endpoint, status).Inc()
}

// RecordExternalAPIDuration records the duration of an external API call
func (m *PrometheusMetrics) RecordExternalAPIDuration(service, endpoint string, duration time.Duration) {
	m.ExternalAPIDuration.WithLabelValues(service, endpoint).Observe(duration.Seconds())
}
