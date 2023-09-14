package kafka

import (
	"context"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/config"
	tracing "github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/telemetry/tracing/kafka"
	"github.com/hanapedia/the-bench/service-unit/stateless/pkg/utils"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/segmentio/kafka-go"
)

type KafkaProducerAdapter struct {
	Writer *kafka.Writer
}

func (kpa KafkaProducerAdapter) Call(ctx context.Context) (string, error) {
	// prepare payload
	payload, err := utils.GenerateRandomString(constants.PayloadSize)
	if err != nil {
		return "", err
	}
	message := kafka.Message{
		Value: []byte(payload),
	}

	// add trace context if tracing is enabled
	if config.GetEnvs().TRACING {
		span := tracing.CreateKafkaProducerSpan(ctx, message)
		defer span.End()
	}

	if err = kpa.Writer.WriteMessages(ctx, message); err != nil {
		return "Failed to produce", err
	}
	return "Successfully produced", nil
}
