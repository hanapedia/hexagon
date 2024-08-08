package kafka

import (
	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/object/crd"
	"github.com/hanapedia/hexagon/pkg/operator/object/factory"
)

// CreateKafkaTopic creates kafka topic
func CreateKafkaTopic(topic string, cc *v1.ClusterConfig) *crd.KafkaTopic {
	kafkaTopicArgs := factory.KafkaTopicArgs{
		Topic:       topic,
		Namespace:   cc.Kafka.Namespace,
		ClusterName: cc.Kafka.ClusterName,
		Partitions:  cc.Kafka.Partitions,
		Replicas:    cc.Kafka.Replications,
	}
	kafkaTopic := factory.NewKafkaTopic(&kafkaTopicArgs)
	return &kafkaTopic
}
