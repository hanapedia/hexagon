package kafka

import (
	"context"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	tracing "github.com/hanapedia/the-bench/internal/service-unit/infrastructure/telemetry/tracing/kafka"
	"github.com/hanapedia/the-bench/pkg/service-unit/utils"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/segmentio/kafka-go"
)

type KafkaProducerAdapter struct {
	writer *kafka.Writer
	payload constants.PayloadSizeVariant
	ports.SecondaryPortBase
}

func (kpa *KafkaProducerAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	// prepare payload
	payload, err := utils.GeneratePayload(kpa.payload)
	if err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
	}
	message := kafka.Message{
		Value: []byte(payload),
	}

	// add trace context if tracing is enabled
	if config.GetEnvs().TRACING {
		span := tracing.CreateKafkaProducerSpan(ctx, message)
		defer span.End()
	}

	if err = kpa.writer.WriteMessages(ctx, message); err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: &payload,
		Error: nil,
	}
}
