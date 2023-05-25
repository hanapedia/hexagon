package generate

import (
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/object/broker"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/yaml"
)

// GenerateBrokerManifests generates manifest file for kafka topic
func (mg ManifestGenerator) GenerateBrokerManifests(config model.BrokerIngressAdapterConfig) ManifestErrors {
	// Open the manifestFile in append mode and with write-only permissions
	manifestFile, err := createFile(mg.Output)
	if err != nil {
		return ManifestErrors{
			broker: []BrokerManifestError{
				NewBrokerManifestError(config, "Unable to open output file."),
			},
		}
	}
	defer manifestFile.Close()

	deployment := broker.CreateKafkaTopic(config.Topic)
	deploymentYAML := yaml.GenerateManifest(deployment)
	_, err = manifestFile.WriteString(formatManifest(deploymentYAML))
	if err != nil {
		return ManifestErrors{
			broker: []BrokerManifestError{
				NewBrokerManifestError(config, "Failed to write deployment manifest"),
			},
		}
	}
	return ManifestErrors{}
}
