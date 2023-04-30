package kafka

import (
	"context"
	"log"
	"reflect"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/common"
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
		log.Printf("message at offset %d: %s = %s", message.Offset, string(message.Key), string(message.Value))
		egressAdapterErrors := common.TaskSetHandler(kca.kafkaConsumer.handler.TaskSets)
		for _, egressAdapterError := range egressAdapterErrors {
			log.Printf("Invocating %s failed: %s",
				reflect.TypeOf(egressAdapterError.EgressAdapter).Elem().Name(),
				egressAdapterError.Error,
			)
		}
	}
	return err
}

func (kca KafkaConsumerAdapter) Register(handler *core.IngressAdapterHandler) error {
	kca.kafkaConsumer.handler = handler
	return nil
}
