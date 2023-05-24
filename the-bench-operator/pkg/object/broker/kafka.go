package broker

import (
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/object/factory"
)

// CreateKafkaTopic creates kafka topic
func CreateKafkaTopic(topic string) *factory.KafkaTopic {
	kafkaTopicArgs := factory.KafkaTopicArgs{
		Topic:       topic,
		Namespace:   factory.KAFKA_NAMESPACE,
		ClusterName: factory.KAFKA_CLUSTER_NAME,
		Replicas:    factory.KAFKA_REPLICATIONS,
		Partitions:  factory.KAFKA_PARTITIONS,
	}
	kafkaTopic := factory.KafkaTopicFactory(&kafkaTopicArgs)
	return &kafkaTopic
}
