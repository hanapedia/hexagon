package templates

type KafkaTopicTemplateArgs struct {
	Topic       string
	ClusterName string
	Namespace   string
	Partitions  int
	Replicas    int
}

const KafkaTopicTemplate = `---
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: {{ .Topic }}
  labels:
    strimzi.io/cluster: {{ .ClusterName }}
  namespace: {{ .Namespace }}
spec:
  partitions: {{ .Partitions }}
  replicas: {{ .Replicas }}
`
