package serviceunit

import (
	"os"

	"github.com/hanapedia/hexagon/internal/hexctl/manifest/core"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/operator/manifest/serviceunit"
	"github.com/hanapedia/hexagon/pkg/operator/yaml"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type ServiceUnitManifest struct {
	deployment     *appsv1.Deployment
	service        *corev1.Service
	configMap      *corev1.ConfigMap
}

func NewServiceUnitManifest(suc *model.ServiceUnitConfig, cc *model.ClusterConfig, configPath string) *ServiceUnitManifest {
	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Logger.Panic("Failed to read config file.")
	}
	manifest := ServiceUnitManifest{
		deployment:     serviceunit.CreateStatelessUnitDeployment(suc, cc),
		service:        serviceunit.CreateStatelessUnitService(suc, cc),
		configMap:      serviceunit.CreateStatelessUnitYamlConfigMap(suc, cc, string(data)),
	}

	return &manifest
}

func (sum *ServiceUnitManifest) Generate(config *model.ServiceUnitConfig, path string) core.ManifestErrors {
	// Open the manifestFile in append mode and with write-only permissions
	file, err := core.CreateFile(path)
	if err != nil {
		return core.ManifestErrors{
			Stateless: []core.StatelessManifestError{
				core.NewStatelessManifestError(config, "Unable to open output file."),
			},
		}
	}
	defer file.Close()

	deploymentYaml := yaml.GenerateManifest(sum.deployment)
	_, err = file.WriteString(core.FormatManifest(deploymentYaml))
	if err != nil {
		return core.ManifestErrors{
			Stateless: []core.StatelessManifestError{
				core.NewStatelessManifestError(config, "Failed to write deployment manifest"),
			},
		}
	}

	serviceYaml := yaml.GenerateManifest(sum.service)
	_, err = file.WriteString(core.FormatManifest(serviceYaml))
	if err != nil {
		return core.ManifestErrors{
			Stateless: []core.StatelessManifestError{
				core.NewStatelessManifestError(config, "Failed to write service manifest"),
			},
		}
	}

	configMapYaml := yaml.GenerateManifest(sum.configMap)
	_, err = file.WriteString(core.FormatManifest(configMapYaml))
	if err != nil {
		return core.ManifestErrors{
			Stateless: []core.StatelessManifestError{
				core.NewStatelessManifestError(config, "Failed to write config map manifest"),
			},
		}
	}
	return core.ManifestErrors{}
}
