package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/core/runtime"
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/primary"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	tracing "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing/kafka"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumerAdapter struct {
	addr          string
	kafkaConsumer *KafkaConsumer
}

type KafkaConsumer struct {
	reader  *kafka.Reader
	handler *primary.PrimaryHandler
}

func NewKafkaConsumerAdapter(topic, group string) *KafkaConsumerAdapter {
	kafkaConsumer := NewKafkaConsumer(topic, group)
	adapter := KafkaConsumerAdapter{addr: config.GetKafkaBrokerAddr(), kafkaConsumer: kafkaConsumer}
	return &adapter
}

func NewKafkaConsumer(topic, group string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{config.GetKafkaBrokerAddr()},
		GroupID:     group,
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
	})
	return &KafkaConsumer{reader: reader}
}

func (kca KafkaConsumerAdapter) Serve(ctx context.Context, wg *sync.WaitGroup) error {
	var err error
ConsumerLoop:
	for {
		message, err := kca.kafkaConsumer.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				logger.Logger.Infof("Context cancelled, Kafka Consumer shutting.")
				kca.kafkaConsumer.reader.Close()
				wg.Done()
			}
			break ConsumerLoop
		}
		startTime := time.Now()

		ctx := context.Background()

		// propagate trace header if tracing is enabled
		ctx, span := tracing.CreateKafkaConsumerSpan(ctx, &message)

		// call tasks
		_ = runtime.TaskSetHandler(ctx, kca.kafkaConsumer.handler)
		// TODO: error handle consumer error
		/* if result.ShouldFail { */
		/* } */

		kca.log(ctx, startTime)
		if span != nil {
			(*span).End()
		}
	}
	return err
}

func (kca *KafkaConsumerAdapter) Register(handler *primary.PrimaryHandler) error {
	kca.kafkaConsumer.handler = handler
	return nil
}

func (kca *KafkaConsumerAdapter) log(ctx context.Context, startTime time.Time) {
	elapsed := time.Since(startTime).Milliseconds()
	logger.Logger.WithContext(ctx).Infof("Message consumed | %-30s | %10v ms", kca.kafkaConsumer.handler.GetId(), elapsed)
}
