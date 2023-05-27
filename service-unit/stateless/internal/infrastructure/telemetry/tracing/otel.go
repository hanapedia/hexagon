package tracing

import (
	"context"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracer(name string) {
	ctx := context.Background()

	// Set up a trace provider
	traceExporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(),
	)
	if err != nil {
		logger.Logger.Fatalf("failed to create trace exporter: %v", err)
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(name),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(1)),
	)

	otel.SetTracerProvider(tracerProvider)

	// Set the global text map propagator to tracecontext.
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
