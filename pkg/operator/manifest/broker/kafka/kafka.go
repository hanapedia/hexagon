package kafka

import (
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/object/crd"
	"github.com/hanapedia/the-bench/pkg/operator/object/factory"
)

// CreateKafkaTopic creates kafka topic
func CreateKafkaTopic(topic string) *crd.KafkaTopic {
	kafkaTopicArgs := factory.KafkaTopicArgs{
		Topic:       topic,
		Namespace:   model.KAFKA_NAMESPACE,
		ClusterName: model.KAFKA_CLUSTER_NAME,
		Replicas:    model.KAFKA_REPLICATIONS,
		Partitions:  model.KAFKA_PARTITIONS,
	}
	kafkaTopic := factory.NewKafkaTopic(&kafkaTopicArgs)
	return &kafkaTopic
}
