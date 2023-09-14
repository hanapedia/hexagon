package kafka

import (
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/segmentio/kafka-go"
)

type KafkaProducerClient struct {
	Client *kafka.Writer
}

// Only cosumer client is included for Kafka
func NewKafkaClient(addr string, topic string) ports.SecondaryAdapter {
	client := kafka.Writer{
		Addr:     kafka.TCP(addr),
		Topic:    topic,
		Balancer: &kafka.RoundRobin{},
	}
	return KafkaProducerClient{Client: &client}
}

func (kafkaProducerClient KafkaProducerClient) Close() {
	kafkaProducerClient.Client.Close()
}
