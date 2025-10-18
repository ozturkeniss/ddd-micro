package monitoring

import (
	"github.com/google/wire"
)

// ProviderSet is a provider set for monitoring infrastructure
var ProviderSet = wire.NewSet(
	NewPrometheusMetrics,
	NewJaegerTracer,
)
