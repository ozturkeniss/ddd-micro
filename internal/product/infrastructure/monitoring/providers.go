package monitoring

import (
	"github.com/google/wire"
)

// ProviderSet is a provider set for monitoring infrastructure
var ProviderSet = wire.NewSet(
	NewPrometheusMetrics,
	ProvideJaegerTracer,
)

// ProvideJaegerTracer provides Jaeger tracer for product service
func ProvideJaegerTracer() (*JaegerTracer, error) {
	return NewJaegerTracer("product-service")
}
