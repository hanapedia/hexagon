package kafka

import (
	"context"
	"log"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumerAdapter struct {
	kafkaConsumer *KafkaConsumer
}

type KafkaConsumer struct {
	addr   string
	reader *kafka.Reader
}

func NewKafkaConsumerAdapter() KafkaConsumerAdapter {
	return KafkaConsumerAdapter{kafkaConsumer: NewKafkaConsumer()}
}

func NewKafkaConsumer() *KafkaConsumer {
	return &KafkaConsumer{addr: constants.KafkaBrokerAddr}
}

func (kca KafkaConsumerAdapter) Serve() error {
	var err error
	for {
		message, err := kca.kafkaConsumer.reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		log.Printf("message at offset %d: %s = %s", message.Offset, string(message.Key), string(message.Value))
	}
	return err
}

func (kca KafkaConsumerAdapter) Register(handler *core.Handler) error {
	kca.kafkaConsumer.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kca.kafkaConsumer.addr},
		Topic:       handler.Action,
		StartOffset: kafka.FirstOffset,
	})
	return nil
}
