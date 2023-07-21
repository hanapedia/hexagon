package kafka

import (
	"context"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/common"
	tracing "github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/telemetry/tracing/kafka"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumerAdapter struct {
	addr          string
	kafkaConsumer *KafkaConsumer
}

type KafkaConsumer struct {
	reader  *kafka.Reader
	handler *core.IngressAdapterHandler
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
		egressAdapterErrors := common.TaskSetHandler(ctx, kca.kafkaConsumer.handler.TaskSets)
		core.LogEgressAdapterErrors(&egressAdapterErrors)
	}
	return err
}

func (kca KafkaConsumerAdapter) Register(serviceName string, handler *core.IngressAdapterHandler) error {
	kca.kafkaConsumer.handler = handler
	return nil
}
