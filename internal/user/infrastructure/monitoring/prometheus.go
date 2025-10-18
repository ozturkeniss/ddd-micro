package monitoring

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics holds all the prometheus metrics for user service
type PrometheusMetrics struct {
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPRequestsInFlight  *prometheus.GaugeVec
	ActiveUsers           prometheus.Gauge
	UserRegistrations     prometheus.Counter
	UserLogins            prometheus.Counter
	UserLoginFailures     prometheus.Counter
	DatabaseConnections   prometheus.Gauge
	DatabaseQueryDuration *prometheus.HistogramVec
}

// NewPrometheusMetrics creates a new instance of PrometheusMetrics
func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "user_service_http_requests_total",
				Help: "Total number of HTTP requests to user service",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "user_service_http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint", "status_code"},
		),
		HTTPRequestsInFlight: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "user_service_http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
			[]string{"method", "endpoint"},
		),
		ActiveUsers: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "user_service_active_users",
				Help: "Current number of active users",
			},
		),
		UserRegistrations: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "user_service_registrations_total",
				Help: "Total number of user registrations",
			},
		),
		UserLogins: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "user_service_logins_total",
				Help: "Total number of successful user logins",
			},
		),
		UserLoginFailures: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "user_service_login_failures_total",
				Help: "Total number of failed user logins",
			},
		),
		DatabaseConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "user_service_database_connections_active",
				Help: "Current number of active database connections",
			},
		),
		DatabaseQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "user_service_database_query_duration_seconds",
				Help:    "Duration of database queries in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
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

// RecordUserRegistration increments the user registration counter
func (m *PrometheusMetrics) RecordUserRegistration() {
	m.UserRegistrations.Inc()
}

// RecordUserLogin increments the user login counter
func (m *PrometheusMetrics) RecordUserLogin() {
	m.UserLogins.Inc()
}

// RecordUserLoginFailure increments the user login failure counter
func (m *PrometheusMetrics) RecordUserLoginFailure() {
	m.UserLoginFailures.Inc()
}

// SetActiveUsers sets the active users gauge
func (m *PrometheusMetrics) SetActiveUsers(count float64) {
	m.ActiveUsers.Set(count)
}

// SetDatabaseConnections sets the database connections gauge
func (m *PrometheusMetrics) SetDatabaseConnections(count float64) {
	m.DatabaseConnections.Set(count)
}

// RecordDatabaseQuery records the duration of a database query
func (m *PrometheusMetrics) RecordDatabaseQuery(operation, table string, duration time.Duration) {
	m.DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}
