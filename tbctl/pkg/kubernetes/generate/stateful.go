package generate

import (
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/object/stateful"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/yaml"
)

// GenerateStatefulManifests creates kubernetes manifest for stateful service units
// currently only Mongo is supported
func (mg ManifestGenerator) GenerateStatefulManifests() ManifestErrors {
	// Open the manifestFile in append mode and with write-only permissions
	outPath := mg.getFilePath(mg.ServiceUnitConfig.Name, "stateful")
	manifestFile, err := createFile(outPath)
	if err != nil {
		return ManifestErrors{
			stateful: []StatefulManifestError{
				NewStatefulManifestError(mg.ServiceUnitConfig, "Unable to open output file."),
			},
		}
	}
	defer manifestFile.Close()

	deployment := stateful.CreateMongoDeployment(mg.ServiceUnitConfig.Name, mg.ServiceUnitConfig.Version)
	deploymentYAML := yaml.GenerateManifest(deployment)
	_, err = manifestFile.WriteString(formatManifest(deploymentYAML))
	if err != nil {
		return ManifestErrors{
			stateful: []StatefulManifestError{
				NewStatefulManifestError(mg.ServiceUnitConfig, "Failed to write deployment manifest"),
			},
		}
	}

	service := stateful.CreateMongoService(mg.ServiceUnitConfig.Name)
	serviceYAML := yaml.GenerateManifest(service)
	_, err = manifestFile.WriteString(formatManifest(serviceYAML))
	if err != nil {
		return ManifestErrors{
			stateful: []StatefulManifestError{
				NewStatefulManifestError(mg.ServiceUnitConfig, "Failed to write service manifest"),
			},
		}
	}

	configMap := stateful.CreateMongoEnvConfigMap(mg.ServiceUnitConfig.Name, MongoEnvs)
	configMapYAML := yaml.GenerateManifest(configMap)
	_, err = manifestFile.WriteString(formatManifest(configMapYAML))
	if err != nil {
		return ManifestErrors{
			stateful: []StatefulManifestError{
				NewStatefulManifestError(mg.ServiceUnitConfig, "Failed to write configmap manifest"),
			},
		}
	}
	return ManifestErrors{}
}
