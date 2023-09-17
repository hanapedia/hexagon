package kafka

import (
	"github.com/segmentio/kafka-go"
)

type KafkaProducerClient struct {
	Client *kafka.Writer
}

// Only cosumer client is included for Kafka
func NewKafkaClient(addr string, topic string) KafkaProducerClient {
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