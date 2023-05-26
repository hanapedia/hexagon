package generate

import (
	"encoding/json"
	"errors"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/object/load_generator"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/yaml"
)

// GenerateLoadGeneratorManifests generates loadgenerator manifest
func (mg ManifestGenerator) GenerateLoadGeneratorManifests() ManifestErrors {
	// Open the manifestFile in append mode and with write-only permissions
	outPath := mg.getFilePath(mg.ServiceUnitConfig.Name, "load-generator")
	manifestFile, err := createFile(outPath)
	if err != nil {
		return ManifestErrors{
			loadGenerator: []LoadGeneratorManifestError{
				NewLoadGeneratorManifestError(mg.ServiceUnitConfig, "Unable to open output file."),
			},
		}
	}
	defer manifestFile.Close()

	deployment := loadgenerator.CreateLoadGeneratorDeployment(mg.ServiceUnitConfig.Name)
	deploymentYAML := yaml.GenerateManifest(deployment)
	_, err = manifestFile.WriteString(formatManifest(deploymentYAML))
	if err != nil {
		return ManifestErrors{
			loadGenerator: []LoadGeneratorManifestError{
				NewLoadGeneratorManifestError(mg.ServiceUnitConfig, "Failed to write deployment manifest"),
			},
		}
	}

	service := loadgenerator.CreateLoadGeneratorService(mg.ServiceUnitConfig.Name)
	serviceYAML := yaml.GenerateManifest(service)
	_, err = manifestFile.WriteString(formatManifest(serviceYAML))
	if err != nil {
		return ManifestErrors{
			loadGenerator: []LoadGeneratorManifestError{
				NewLoadGeneratorManifestError(mg.ServiceUnitConfig, "Failed to write service manifest"),
			},
		}
	}

	rawConfig, err := mg.GenerateConfigJson()
	if err != nil {
		return ManifestErrors{
			loadGenerator: []LoadGeneratorManifestError{
				NewLoadGeneratorManifestError(mg.ServiceUnitConfig, "Failed to create config json"),
			},
		}
	}

	rawRoutes, err := mg.GenerateRoutesJson()
	if err != nil {
		return ManifestErrors{
			loadGenerator: []LoadGeneratorManifestError{
				NewLoadGeneratorManifestError(mg.ServiceUnitConfig, "Failed to create routes json"),
			},
		}
	}

	configMap := loadgenerator.CreateLoadGeneratorYamlConfigMap(
		mg.ServiceUnitConfig.Name,
		string(rawConfig),
		string(rawRoutes),
	)
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

// GenerateConfigJson generate json content for config
func (mg ManifestGenerator) GenerateConfigJson() ([]byte, error) {
	if mg.ServiceUnitConfig.Gateway == nil {
		return nil, errors.New("Gateway config not found")
	}
	config := loadgenerator.CreateLoadGeneratorConfig(
		mg.ServiceUnitConfig.Gateway.VirtualUsers,
		mg.ServiceUnitConfig.Gateway.Duration,
		mg.ServiceUnitConfig.Name,
	)
	return json.Marshal(config)
}

// GenerateConfigJson generate json content for config
func (mg ManifestGenerator) GenerateRoutesJson() ([]byte, error) {
	if mg.ServiceUnitConfig.Gateway == nil {
		return nil, errors.New("Gateway config not found")
	}
	var routes []loadgenerator.Route
	for _, ingressAdapter := range mg.ServiceUnitConfig.IngressAdapterConfigs {
		if ingressAdapter.StatelessIngressAdapterConfig != nil {
			var weight int32 = 1
			if ingressAdapter.StatelessIngressAdapterConfig.Weight != nil {
				weight = *ingressAdapter.StatelessIngressAdapterConfig.Weight
			}
			route := loadgenerator.CreateLoadGeneratorRoutes(
				ingressAdapter.StatelessIngressAdapterConfig.Route,
				constants.GetHttpMethodFromAction(ingressAdapter.StatelessIngressAdapterConfig.Action),
				weight,
			)
			routes = append(routes, route)
		}
	}
	return json.Marshal(routes)
}
