package kafka

import (
	"context"
	"fmt"

	tracing "github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/telemetry/tracing"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// CreateKafkaConsumerSpan creates consumer span
func CreateKafkaConsumerSpan(message kafka.Message) trace.Span {
	// Extract the span context from the message headers.
	propagator := propagation.TraceContext{}
	carrier := KafkaCarrier(message.Headers)
	extractedCtx := propagator.Extract(context.Background(), carrier)

	// Start a new span for the Kafka message consumption.
	tracer := tracing.GetTracer()
	_, consumerSpan := tracer.Start(extractedCtx, fmt.Sprintf("%s receive", message.Topic))

	consumerSpan.SetAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination_kind", "topic"),
		attribute.String("messaging.operation", "receive"),
		attribute.String("messaging.destination", message.Topic),
		attribute.Int("messaging.kafka.partition", message.Partition),
		attribute.Int64("messaging.kafka.offset", message.Offset),
	)
	return consumerSpan
}

// CreateKafkaProducerSpan creates consumer span
func CreateKafkaProducerSpan(ctx context.Context, message kafka.Message) trace.Span {
	// prepare span
	tracer := tracing.GetTracer()
	_, producerSpan := tracer.Start(ctx, fmt.Sprintf("%s publish", message.Topic))

	// create span context
	spanCtx := trace.ContextWithSpan(ctx, producerSpan)

	// add trace context
	propagator := propagation.TraceContext{}
	carrier := KafkaCarrier(message.Headers)
	propagator.Inject(spanCtx, carrier)

	producerSpan.SetAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination_kind", "topic"),
		attribute.String("messaging.operation", "publish"),
		attribute.String("messaging.destination", message.Topic),
	)

	return producerSpan
}
