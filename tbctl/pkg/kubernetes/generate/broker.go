package generate

import (
	"fmt"

	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/tbctl/pkg/kubernetes/templates"
)

func GenerateBrokerManifests(dir string, serviceUnitConfig model.ServiceUnitConfig) ManifestErrors {
	var brokerManifestErrors []BrokerManifestError
	for _, ingressAdapter := range serviceUnitConfig.IngressAdapterConfigs {
		if ingressAdapter.BrokerIngressAdapterConfig == nil {
			continue
		}
		// generate kafka topic manifest
		if ingressAdapter.BrokerIngressAdapterConfig.Variant == constants.KAFKA {
			err := generateKafkaTopicManifest(dir, *ingressAdapter.BrokerIngressAdapterConfig)
			if err != nil {
				brokerManifestErrors = append(
					brokerManifestErrors,
					NewBrokerManifestError(*ingressAdapter.BrokerIngressAdapterConfig, err.Error()),
				)
			}
			continue
		}
	}
	return ManifestErrors{broker: brokerManifestErrors}
}

func generateKafkaTopicManifest(dir string, brokerAdapterConfig model.BrokerAdapterConfig) error {
	kafkaTopicTemplateArgs := templates.KafkaTopicTemplateArgs{
		Topic:       brokerAdapterConfig.Topic,
		ClusterName: KAFKA_CLUSTER_NAME,
		Namespace:   KAFKA_NAMESPACE,
		Partitions:  KAFKA_PARTITIONS,
		Replicas:    KAFKA_REPLICATIONS,
	}
	err := RenderAndSave(
		dir,
		fmt.Sprintf("%s-kafka-topic", brokerAdapterConfig.Topic),
		templates.KafkaTopicTemplate,
		kafkaTopicTemplateArgs,
	)
	return err
}
