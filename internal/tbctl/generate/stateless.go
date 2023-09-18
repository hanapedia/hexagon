package generate

import (
	"os"

	"github.com/hanapedia/the-bench/pkg/operator/object/usecases"
	"github.com/hanapedia/the-bench/pkg/operator/yaml"
)

// GenerateStatelessManifests generates stateless manifest
func (mg ManifestGenerator) GenerateStatelessManifests() ManifestErrors {
	// Open the manifestFile in append mode and with write-only permissions
	outPath := mg.getFilePath(mg.ServiceUnitConfig.Name, "stateless")
	manifestFile, err := createFile(outPath)
	if err != nil {
		return ManifestErrors{
			stateless: []StatelessManifestError{
				NewStatelessManifestError(mg.ServiceUnitConfig, "Unable to open output file."),
			},
		}
	}
	defer manifestFile.Close()

	deployment := usecases.CreateStatelessUnitDeployment(mg.ServiceUnitConfig.Name, mg.ServiceUnitConfig.Version)
	deploymentYAML := yaml.GenerateManifest(deployment)
	_, err = manifestFile.WriteString(formatManifest(deploymentYAML))
	if err != nil {
		return ManifestErrors{
			stateless: []StatelessManifestError{
				NewStatelessManifestError(mg.ServiceUnitConfig, "Failed to write deployment manifest"),
			},
		}
	}

	service := usecases.CreateStatelessUnitService(mg.ServiceUnitConfig.Name)
	serviceYAML := yaml.GenerateManifest(service)
	_, err = manifestFile.WriteString(formatManifest(serviceYAML))
	if err != nil {
		return ManifestErrors{
			stateless: []StatelessManifestError{
				NewStatelessManifestError(mg.ServiceUnitConfig, "Failed to write service manifest"),
			},
		}
	}

	data, err := os.ReadFile(mg.Input)
	if err != nil {
		return ManifestErrors{
			stateless: []StatelessManifestError{
				NewStatelessManifestError(mg.ServiceUnitConfig, "Unable to read config file."),
			},
		}
	}
	configMap := usecases.CreateStatelessUnitYamlConfigMap(mg.ServiceUnitConfig.Name, string(data))
	configMapYAML := yaml.GenerateManifest(configMap)
	_, err = manifestFile.WriteString(formatManifest(configMapYAML))
	if err != nil {
		return ManifestErrors{
			stateless: []StatelessManifestError{
				NewStatelessManifestError(mg.ServiceUnitConfig, "Failed to write configmap manifest"),
			},
		}
	}
	return ManifestErrors{}
}
