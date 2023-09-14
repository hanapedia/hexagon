package kafka

import (
	"context"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/core/runtime"
	tracing "github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/telemetry/tracing/kafka"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumerAdapter struct {
	addr          string
	kafkaConsumer *KafkaConsumer
}

type KafkaConsumer struct {
	reader  *kafka.Reader
	handler *ports.PrimaryHandler
}

func NewKafkaConsumerAdapter(topic string) KafkaConsumerAdapter {
	kafkaConsumer := NewKafkaConsumer(topic)
	return KafkaConsumerAdapter{addr: config.GetKafkaBrokerAddr(), kafkaConsumer: kafkaConsumer}
}

func NewKafkaConsumer(topic string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{config.GetKafkaBrokerAddr()},
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
	})
	return &KafkaConsumer{reader: reader}
}

func (kca KafkaConsumerAdapter) Serve() error {
	var err error
	for {
		message, err := kca.kafkaConsumer.reader.ReadMessage(context.Background())
		if err != nil {
			break
		}

		ctx := context.Background()

		// propagate trace header if tracing is enabled
		if config.GetEnvs().TRACING {
			span := tracing.CreateKafkaConsumerSpan(message)
			defer span.End()
		}

		// call tasks
		secondaryAdapterErrors := runtime.TaskSetHandler(ctx, kca.kafkaConsumer.handler.TaskSets)
		ports.LogSecondaryPortErrors(&secondaryAdapterErrors)
	}
	return err
}

func (kca KafkaConsumerAdapter) Register(serviceName string, handler *ports.PrimaryHandler) error {
	kca.kafkaConsumer.handler = handler
	return nil
}
