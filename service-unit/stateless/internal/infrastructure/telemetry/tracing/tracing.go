package tracing

import (
	"context"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitTracer(name, collectorUrl string) *sdktrace.TracerProvider {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, collectorUrl,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Logger.Fatalf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
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

	// return provider so that it can be shutdown
	return tracerProvider
}

// get the default tracer for creating original spans
func GetTracer() trace.Tracer {
	return otel.Tracer("the-bench")
}
