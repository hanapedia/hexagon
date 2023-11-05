package kafka

import (
	"context"
	"fmt"

	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	tracing "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// CreateKafkaConsumerSpan creates consumer span
func CreateKafkaConsumerSpan(ctx context.Context, message *kafka.Message) (context.Context, *trace.Span) {
	// return early if tracing is disabled
	if !config.GetEnvs().TRACING {
		return ctx, nil
	}

	// Extract the span context from the message headers.
	propagator := propagation.TraceContext{}
	carrier := &KafkaCarrier{Headers: message.Headers}
	extractedCtx := propagator.Extract(ctx, carrier)

	// Start a new span for the Kafka message consumption.
	tracer := tracing.GetTracer()
	tracedCtx, consumerSpan := tracer.Start(extractedCtx, fmt.Sprintf("%s receive", message.Topic))

	consumerSpan.SetAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination_kind", "topic"),
		attribute.String("messaging.operation", "receive"),
		attribute.String("messaging.destination", message.Topic),
		attribute.Int("messaging.kafka.partition", message.Partition),
		attribute.Int64("messaging.kafka.offset", message.Offset),
	)
	return tracedCtx, &consumerSpan
}

// CreateKafkaProducerSpan creates producer span
func CreateKafkaProducerSpan(ctx context.Context, message *kafka.Message) (context.Context, *trace.Span) {
	// return early if tracing is disabled
	if !config.GetEnvs().TRACING {
		return ctx, nil
	}

	// prepare span
	tracer := tracing.GetTracer()
	_, producerSpan := tracer.Start(ctx, fmt.Sprintf("%s publish", message.Topic))

	// create span context
	spanCtx := trace.ContextWithSpan(ctx, producerSpan)

	// add trace context
	propagator := propagation.TraceContext{}
	carrier := &KafkaCarrier{Headers: message.Headers}
	propagator.Inject(spanCtx, carrier)
	message.Headers = []kafka.Header(carrier.Headers)

	producerSpan.SetAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination_kind", "topic"),
		attribute.String("messaging.operation", "publish"),
		attribute.String("messaging.destination", message.Topic),
	)

	return spanCtx, &producerSpan
}
