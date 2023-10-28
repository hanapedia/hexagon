package loadgenerator

import (
	"github.com/hanapedia/hexagon/internal/hexctl/manifest/core"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/operator/manifest/loadgenerator/k6"
	"github.com/hanapedia/hexagon/pkg/operator/yaml"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type LoadGeneratorManifest struct {
	deployment *appsv1.Deployment
	service    *corev1.Service
	configMap  *corev1.ConfigMap
}

func NewLoadGeneratorManifest(config *model.ServiceUnitConfig) *LoadGeneratorManifest {
	if !core.HasGatewayConfig(config) {
		logger.Logger.Panic("Gateway config not found.")
	}

	manifest := LoadGeneratorManifest{
		deployment: k6.CreateLoadGeneratorDeployment(config),
		service:    k6.CreateLoadGeneratorService(config),
		configMap:  k6.CreateLoadGeneratorYamlConfigMap(config),
	}

	return &manifest
}

func (sum *LoadGeneratorManifest) Generate(config *model.ServiceUnitConfig, path string) core.ManifestErrors {
	// Open the manifestFile in append mode and with write-only permissions
	file, err := core.CreateFile(path)
	if err != nil {
		return core.ManifestErrors{
			LoadGenerator: []core.LoadGeneratorManifestError{
				core.NewLoadGeneratorManifestError(config, "Unable to open output file."),
			},
		}
	}
	defer file.Close()

	deploymentYaml := yaml.GenerateManifest(sum.deployment)
	_, err = file.WriteString(core.FormatManifest(deploymentYaml))
	if err != nil {
		return core.ManifestErrors{
			LoadGenerator: []core.LoadGeneratorManifestError{
				core.NewLoadGeneratorManifestError(config, "Failed to write deployment manifest"),
			},
		}
	}

	serviceYaml := yaml.GenerateManifest(sum.service)
	_, err = file.WriteString(core.FormatManifest(serviceYaml))
	if err != nil {
		return core.ManifestErrors{
			LoadGenerator: []core.LoadGeneratorManifestError{
				core.NewLoadGeneratorManifestError(config, "Failed to write service manifest"),
			},
		}
	}

	configMapYaml := yaml.GenerateManifest(sum.configMap)
	_, err = file.WriteString(core.FormatManifest(configMapYaml))
	if err != nil {
		return core.ManifestErrors{
			LoadGenerator: []core.LoadGeneratorManifestError{
				core.NewLoadGeneratorManifestError(config, "Failed to write config map manifest"),
			},
		}
	}

	return core.ManifestErrors{}
}
