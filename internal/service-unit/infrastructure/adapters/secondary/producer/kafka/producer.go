package kafka

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	tracing "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing/kafka"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"github.com/segmentio/kafka-go"
)

type kafkaProducerAdapter struct {
	writer      *kafka.Writer
	payloadSize int64
	secondary.SecondaryPortBase
}

func (kpa *kafkaProducerAdapter) Call(ctx context.Context) secondary.SecondaryPortCallResult {
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
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return secondary.SecondaryPortCallResult{
		Payload: &payload,
		Error:   nil,
	}
}
