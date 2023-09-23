package kafka

import (
	"context"
	"time"

	"github.com/hanapedia/the-bench/internal/service-unit/application/core/runtime"
	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	tracing "github.com/hanapedia/the-bench/internal/service-unit/infrastructure/telemetry/tracing/kafka"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
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

func NewKafkaConsumerAdapter(topic string) *KafkaConsumerAdapter {
	kafkaConsumer := NewKafkaConsumer(topic)
	adapter := KafkaConsumerAdapter{addr: config.GetKafkaBrokerAddr(), kafkaConsumer: kafkaConsumer}
	return &adapter
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
		startTime := time.Now()

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
		errs := runtime.TaskSetHandler(ctx, kca.kafkaConsumer.handler.TaskSet)
		if errs != nil {
			for _, err := range errs {
				kca.kafkaConsumer.handler.LogTaskError(ctx, err)
			}
		}

		kca.log(ctx, time.Since(startTime).Milliseconds())
	}
	return err
}

func (kca KafkaConsumerAdapter) Register(serviceName string, handler *ports.PrimaryHandler) error {
	kca.kafkaConsumer.handler = handler
	return nil
}

func (kca KafkaConsumerAdapter) log(ctx context.Context, elapsed int64) {
	logger.Logger.WithContext(ctx).Infof("Message consumed | %-30s | %10v ms", kca.kafkaConsumer.handler.GetId(), elapsed)
}
