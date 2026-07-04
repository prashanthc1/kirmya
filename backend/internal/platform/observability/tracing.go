package observability

import (
	"context"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const serviceName = "kirmya"

// InitTracing configures an OTLP/HTTP trace exporter when
// OTEL_EXPORTER_OTLP_ENDPOINT is set; otherwise tracing is a no-op. It returns a
// shutdown function that flushes pending spans (safe to call even when disabled).
func InitTracing(ctx context.Context) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		log.Printf("[otel] OTEL_EXPORTER_OTLP_ENDPOINT not set; tracing disabled")
		return func(context.Context) error { return nil }, nil
	}

	// otlptracehttp reads OTEL_EXPORTER_OTLP_ENDPOINT itself; insecure unless the
	// endpoint is https.
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))
	log.Printf("[otel] tracing enabled, exporting to %s", endpoint)
	return tp.Shutdown, nil
}

// WrapHandler instruments the HTTP handler with OpenTelemetry spans. When tracing
// is disabled the global no-op tracer makes this near-zero cost.
func WrapHandler(h http.Handler) http.Handler {
	return otelhttp.NewHandler(h, "http.server")
}
