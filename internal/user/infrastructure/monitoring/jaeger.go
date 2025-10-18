package monitoring

import (
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// JaegerTracer wraps the Jaeger tracer functionality
type JaegerTracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

// NewJaegerTracer creates a new Jaeger tracer instance
func NewJaegerTracer(serviceName string) (*JaegerTracer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}

	// Set the global tracer
	opentracing.SetGlobalTracer(tracer)

	return &JaegerTracer{
		tracer: tracer,
		closer: closer,
	}, nil
}

// Close closes the tracer
func (jt *JaegerTracer) Close() error {
	if jt.closer != nil {
		return jt.closer.Close()
	}
	return nil
}

// StartSpan creates a new span
func (jt *JaegerTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return jt.tracer.StartSpan(operationName, opts...)
}

// StartSpanFromContext creates a new span from context
func (jt *JaegerTracer) StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}

// Inject injects span context into carrier
func (jt *JaegerTracer) Inject(spanContext opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return jt.tracer.Inject(spanContext, format, carrier)
}

// Extract extracts span context from carrier
func (jt *JaegerTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return jt.tracer.Extract(format, carrier)
}

// JaegerMiddleware returns a Gin middleware for Jaeger tracing
func JaegerMiddleware(tracer *JaegerTracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract span context from headers
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			// Log error but continue
		}

		// Start new span
		span, ctx := tracer.StartSpanFromContext(c.Request.Context(), c.Request.Method+" "+c.FullPath(), ext.RPCServerOption(spanCtx))
		defer span.Finish()

		// Set span tags
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		ext.Component.Set(span, "user-service")
		span.SetTag("http.path", c.FullPath())
		span.SetTag("service.name", "user-service")

		// Update context with span
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Set response tags
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		if c.Writer.Status() >= 400 {
			ext.Error.Set(span, true)
			span.SetTag("error", true)
		}
	}
}

// StartSpanFromGinContext creates a span from Gin context
func StartSpanFromGinContext(c *gin.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(c.Request.Context(), operationName, opts...)
}

// SetSpanTags sets common tags for user service spans
func SetSpanTags(span opentracing.Span, tags map[string]interface{}) {
	for key, value := range tags {
		span.SetTag(key, value)
	}
}

// LogSpanEvent logs an event to the span
func LogSpanEvent(span opentracing.Span, event string, fields ...opentracing.LogField) {
	span.LogFields(fields...)
	span.LogEvent(event)
}

// LogSpanError logs an error to the span
func LogSpanError(span opentracing.Span, err error) {
	ext.Error.Set(span, true)
	span.SetTag("error", true)
	span.LogFields(opentracing.LogField{
		Key:   "error",
		Value: err.Error(),
	})
}
