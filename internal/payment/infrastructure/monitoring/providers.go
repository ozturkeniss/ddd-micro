package monitoring

import (
	"github.com/google/wire"
)

// ProviderSet is a provider set for monitoring infrastructure
var ProviderSet = wire.NewSet(
	NewPrometheusMetrics,
	ProvideJaegerTracer,
)

// ProvideJaegerTracer provides Jaeger tracer for payment service
func ProvideJaegerTracer() (*JaegerTracer, error) {
	return NewJaegerTracer("payment-service")
}
