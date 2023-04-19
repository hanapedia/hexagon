package kafka

import (
	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/segmentio/kafka-go"
)

type KafkaWriterConnection struct {
	Connection *kafka.Writer
}

// Only cosumer connection is included for Kafka
func NewKafkaConnection(addr string, topic string) core.EgressConnection {
	connection := kafka.Writer{
		Addr:     kafka.TCP(addr),
		Topic:    topic,
		Balancer: &kafka.RoundRobin{},
	}
	return KafkaWriterConnection{Connection: &connection}
}

func (kafkaWriterConnection KafkaWriterConnection) Close() {
	kafkaWriterConnection.Close()
}
