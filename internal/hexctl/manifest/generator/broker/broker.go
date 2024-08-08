package broker

import (
	"github.com/hanapedia/hexagon/internal/hexctl/manifest/core"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/operator/manifest/broker/kafka"
	"github.com/hanapedia/hexagon/pkg/operator/object/crd"
	"github.com/hanapedia/hexagon/pkg/operator/yaml"
)

type BrokerManifest struct {
	kafkaTopics []*crd.KafkaTopic
	// add other crds if new broker is added
}

func NewBrokerManifest(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *BrokerManifest {
	consumerAdapters := core.GetConsumerAdapters(suc)
	if len(consumerAdapters) == 0 {
		logger.Logger.Panic("No consumer adapter config found.")
	}

	var manifest BrokerManifest
	for _, consumerAdapter := range consumerAdapters {
		switch consumerAdapter.Variant {
		case constants.KAFKA:
			manifest.kafkaTopics = append(manifest.kafkaTopics, kafka.CreateKafkaTopic(consumerAdapter.Topic, cc))
		default:
			logger.Logger.Panic("Invalid broker variant.")
		}
	}

	return &manifest
}

func (sum *BrokerManifest) Generate(config *model.ServiceUnitConfig, path string) core.ManifestErrors {
	// Open the manifestFile in append mode and with write-only permissions
	file, err := core.CreateFile(path)
	if err != nil {
		return core.ManifestErrors{
			Broker: []core.BrokerManifestError{
				core.NewBrokerManifestFileError(config, "Unable to open output file."),
			},
		}
	}
	defer file.Close()

	for _, kafkaTopic := range sum.kafkaTopics {
		kafkaTopicYaml := yaml.GenerateManifest(kafkaTopic)
		_, err = file.WriteString(core.FormatManifest(kafkaTopicYaml))
		if err != nil {
			return core.ManifestErrors{
				Stateful: []core.StatefulManifestError{
					core.NewStatefulManifestError(config, "Failed to write kafka topic manifest"),
				},
			}
		}
	}

	return core.ManifestErrors{}
}
