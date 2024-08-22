package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/core/runtime"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/consumer"
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
	handler *domain.PrimaryAdapterHandler
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

func (kca KafkaConsumerAdapter) Serve(ctx context.Context, shutdownWg, readyWg *sync.WaitGroup) error {
	var err error
	// Mark Ready
	readyWg.Done()
ConsumerLoop:
	for {
		consumer.SetConsumerAdapterInProgress(
			domain.INC,
			kca.kafkaConsumer.handler.ServiceName,
			kca.kafkaConsumer.handler.ConsumerConfig,
		)
		message, err := kca.kafkaConsumer.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				logger.Logger.Infof("Context cancelled, Kafka Consumer shutting.")
				kca.kafkaConsumer.reader.Close()
				shutdownWg.Done()
			}
			break ConsumerLoop
		}
		startTime := time.Now()

		ctx := context.Background()

		// propagate trace header if tracing is enabled
		ctx, span := tracing.CreateKafkaConsumerSpan(ctx, &message)

		// call tasks
		result := runtime.TaskSetHandler(ctx, kca.kafkaConsumer.handler)
		// TODO: error handle consumer error
		/* if result.ShouldFail { */
		/* } */

		// record metrics
		go consumer.ObserveConsumerAdapterDuration(
			time.Since(startTime),
			kca.kafkaConsumer.handler.ServiceName,
			kca.kafkaConsumer.handler.ConsumerConfig,
			result.ShouldFail,
		)

		consumer.Log(kca.kafkaConsumer.handler, startTime)
		if span != nil {
			(*span).End()
		}
		consumer.SetConsumerAdapterInProgress(
			domain.DEC,
			kca.kafkaConsumer.handler.ServiceName,
			kca.kafkaConsumer.handler.ConsumerConfig,
		)
	}
	return err
}

func (kca *KafkaConsumerAdapter) Register(handler *domain.PrimaryAdapterHandler) error {
	kca.kafkaConsumer.handler = handler
	return nil
}
