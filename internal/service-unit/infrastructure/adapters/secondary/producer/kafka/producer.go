package kafka

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	tracing "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing/kafka"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"github.com/segmentio/kafka-go"
)

type kafkaProducerAdapter struct {
	writer      *kafka.Writer
	payloadSize int64
	ports.SecondaryPortBase
}

func (kpa *kafkaProducerAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	// prepare payload
	payload := utils.GenerateRandomString(kpa.payloadSize)
	message := kafka.Message{
		Value: []byte(payload),
	}

	ctx, span := tracing.CreateKafkaProducerSpan(ctx, &message)
	if span != nil {
		defer (*span).End()
	}

	if err := kpa.writer.WriteMessages(ctx, message); err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: &payload,
		Error:   nil,
	}
}
