package factory

import "github.com/hanapedia/the-bench/pkg/operator/object/crd"

type KafkaTopicArgs struct {
	Topic       string
	ClusterName string
	Namespace   string
	Partitions  int32
	Replicas    int32
}

// NewKafkaTopic creates kafka topic custom resource defined by strimzi.io
func NewKafkaTopic(args *KafkaTopicArgs) crd.KafkaTopic {
	return crd.KafkaTopic{
		TypeMeta: NewTypeMeta("KafkaTopic", "kafka.strimzi.io/v1beta2"),
		ObjectMeta: NewObjectMeta(ObjectMetaOptions{
			Name:      args.Topic,
			Namespace: args.Namespace,
			Labels: map[string]string{
				"strimzi.io/cluster": args.ClusterName,
			},
		},
		),
		Spec: crd.KafkaTopicSpec{
			Partitions: args.Partitions,
			Replicas:   args.Replicas,
		},
	}
}
