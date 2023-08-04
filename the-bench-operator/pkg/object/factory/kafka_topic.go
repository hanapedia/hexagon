package factory

import "github.com/hanapedia/the-bench/the-bench-operator/pkg/object/crd"

type KafkaTopicArgs struct {
	Topic       string
	ClusterName string
	Namespace   string
	Partitions  int32
	Replicas    int32
}

// KafkaTopicFactory creates kafka topic custom resource defined by strimzi.io
func KafkaTopicFactory(args *KafkaTopicArgs) crd.KafkaTopic {
	return crd.KafkaTopic{
		TypeMeta: TypeMetaFactory("KafkaTopic", "kafka.strimzi.io/v1beta2"),
		ObjectMeta: ObjectMetaFactory(ObjectMetaOptions{
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
