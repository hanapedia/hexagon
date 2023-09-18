package usecases

import (
	"github.com/hanapedia/the-bench/pkg/operator/object/crd"
	"github.com/hanapedia/the-bench/pkg/operator/object/factory"
)

// CreateKafkaTopic creates kafka topic
func CreateKafkaTopic(topic string) *crd.KafkaTopic {
	kafkaTopicArgs := factory.KafkaTopicArgs{
		Topic:       topic,
		Namespace:   factory.KAFKA_NAMESPACE,
		ClusterName: factory.KAFKA_CLUSTER_NAME,
		Replicas:    factory.KAFKA_REPLICATIONS,
		Partitions:  factory.KAFKA_PARTITIONS,
	}
	kafkaTopic := factory.NewKafkaTopic(&kafkaTopicArgs)
	return &kafkaTopic
}
