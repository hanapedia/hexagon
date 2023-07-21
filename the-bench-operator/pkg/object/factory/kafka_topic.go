package factory

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CustomResource represents a custom Kubernetes resource.
type KafkaTopic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec KafkaTopicSpec `json:"spec,omitempty"`
}

// KafkaTopicSpec represents the spec of the custom resource.
type KafkaTopicSpec struct {
	Partitions int32 `json:"partitions,omitempty"`
	Replicas   int32 `json:"replicas,omitempty"`
}

type KafkaTopicArgs struct {
	Topic       string
	ClusterName string
	Namespace   string
	Partitions  int32
	Replicas    int32
}

// KafkaTopicFactory creates kafka topic custom resource defined by strimzi.io
func KafkaTopicFactory(args *KafkaTopicArgs) KafkaTopic {
	return KafkaTopic{
		TypeMeta: TypeMetaFactory("KafkaTopic", "kafka.strimzi.io/v1beta2"),
		ObjectMeta: ObjectMetaFactory(ObjectMetaOptions{
			Name:      args.Topic,
			Namespace: args.Namespace,
			Labels: map[string]string{
				"strimzi.io/cluster": args.ClusterName,
			},
		},
		),
		Spec: KafkaTopicSpec{
			Partitions: args.Partitions,
			Replicas:   args.Replicas,
		},
	}
}
