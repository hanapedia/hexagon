package crd

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

