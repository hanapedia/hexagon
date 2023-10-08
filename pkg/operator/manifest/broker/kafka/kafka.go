package kafka

import (
	"github.com/hanapedia/the-bench/pkg/api/defaults"
	"github.com/hanapedia/the-bench/pkg/operator/object/crd"
	"github.com/hanapedia/the-bench/pkg/operator/object/factory"
)

// CreateKafkaTopic creates kafka topic
func CreateKafkaTopic(topic string) *crd.KafkaTopic {
	kafkaTopicArgs := factory.KafkaTopicArgs{
		Topic:       topic,
		Namespace:   defaults.KAFKA_NAMESPACE,
		ClusterName: defaults.KAFKA_CLUSTER_NAME,
		Replicas:    defaults.KAFKA_REPLICATIONS,
		Partitions:  defaults.KAFKA_PARTITIONS,
	}
	kafkaTopic := factory.NewKafkaTopic(&kafkaTopicArgs)
	return &kafkaTopic
}
