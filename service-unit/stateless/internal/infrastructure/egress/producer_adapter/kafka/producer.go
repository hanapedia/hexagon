package kafka

import (
	"context"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/service-unit/stateless/pkg/utils"
	"github.com/segmentio/kafka-go"
)

type KafkaProducerAdapter struct {
	Writer *kafka.Writer
}

func (kpa KafkaProducerAdapter) Call() (string, error) {
	payload, err := utils.GenerateRandomString(constants.PayloadSize)
	if err != nil {
		return "", err
	}
	err = kpa.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(payload),
		},
	)
    return "Successfully produced", err
}
