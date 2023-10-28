package statefulunit

import (
	"github.com/hanapedia/hexagon/internal/hexctl/manifest/core"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/operator/manifest/statefulunit/mongo"
	"github.com/hanapedia/hexagon/pkg/operator/manifest/statefulunit/redis"
	"github.com/hanapedia/hexagon/pkg/operator/yaml"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type StatefulUnitManifest struct {
	deployment *appsv1.Deployment
	service    *corev1.Service
}

func NewStatefulUnitManifest(config *model.ServiceUnitConfig) *StatefulUnitManifest {
	repositoryAdapter := core.GetRepositoryAdapter(config)
	if repositoryAdapter == nil {
		logger.Logger.Panic("No repository adapter found.")
	}

	var manifest StatefulUnitManifest
	switch repositoryAdapter.Variant {
	case constants.MONGO:
		manifest = StatefulUnitManifest{
			deployment: mongo.CreateMongoDeployment(config),
			service: mongo.CreateMongoService(config),
		}
	case constants.REDIS:
		manifest = StatefulUnitManifest{
			deployment: redis.CreateRedisDeployment(config),
			service: redis.CreateRedisService(config),
		}
	default:
		logger.Logger.Panic("Invalid repository variant.")
	}

	return &manifest
}

func (sum *StatefulUnitManifest) Generate(config *model.ServiceUnitConfig, path string) core.ManifestErrors {
	// Open the manifestFile in append mode and with write-only permissions
	file, err := core.CreateFile(path)
	if err != nil {
		return core.ManifestErrors{
			Stateful: []core.StatefulManifestError{
				core.NewStatefulManifestError(config, "Unable to open output file."),
			},
		}
	}
	defer file.Close()

	deploymentYaml := yaml.GenerateManifest(sum.deployment)
	_, err = file.WriteString(core.FormatManifest(deploymentYaml))
	if err != nil {
		return core.ManifestErrors{
			Stateful: []core.StatefulManifestError{
				core.NewStatefulManifestError(config, "Failed to write deployment manifest"),
			},
		}
	}

	serviceYaml := yaml.GenerateManifest(sum.service)
	_, err = file.WriteString(core.FormatManifest(serviceYaml))
	if err != nil {
		return core.ManifestErrors{
			Stateful: []core.StatefulManifestError{
				core.NewStatefulManifestError(config, "Failed to write service manifest"),
			},
		}
	}

	return core.ManifestErrors{}
}
