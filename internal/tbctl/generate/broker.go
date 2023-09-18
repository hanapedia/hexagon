package generate

import (
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/object/usecases"
	"github.com/hanapedia/the-bench/pkg/operator/yaml"
)

// GenerateBrokerManifests generates manifest file for kafka topic
func (mg ManifestGenerator) GenerateBrokerManifests(config model.ConsumerConfig) ManifestErrors {
	// Open the manifestFile in append mode and with write-only permissions
	outPath := mg.getFilePath(config.Topic, "kafka-topic")
	manifestFile, err := createFile(outPath)
	if err != nil {
		return ManifestErrors{
			broker: []BrokerManifestError{
				NewBrokerManifestError(config, "Unable to open output file."),
			},
		}
	}
	defer manifestFile.Close()

	deployment := usecases.CreateKafkaTopic(config.Topic)
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
