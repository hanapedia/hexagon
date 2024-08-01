package tracing

import (
	"context"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(collectorUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Logger.Errorf("failed to create gRPC connection to collector: %s, setting TRACING=false", err)
		config.GetEnvs().TRACING = false
		return nil
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		logger.Logger.Errorf("failed to create trace exporter: %v", err)
		config.GetEnvs().TRACING = false
		return nil
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
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// return provider so that it can be shutdown
	return tracerProvider
}

// get the default tracer for creating original spans
func GetTracer() trace.Tracer {
	return otel.Tracer("github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing")
}
