package redis

import (
	"fmt"

	"github.com/hanapedia/hexagon/pkg/api/defaults"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/object/factory"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// CreateRedisDeployment creates deployment for redis
func CreateRedisDeployment(config *model.ServiceUnitConfig) *appsv1.Deployment {
	replica := config.DeploymentSpec.Replicas
	if replica <= 0 {
		replica = 1
	}

	resource := config.DeploymentSpec.Resource
	if resource == nil {
		resource = getDefaultResource()
	}

	envs := config.DeploymentSpec.EnvVar

	deploymentArgs := factory.DeploymentArgs{
		Name:         config.Name,
		Namespace:    defaults.NAMESPACE,
		Annotations:  map[string]string{"rca": "ignore"},
		Image:        fmt.Sprintf("%s:%s", defaults.REDIS_IMAGE_NAME, config.Version),
		Replicas:     replica,
		Resource:     resource,
		Ports:        map[string]int32{"redis": defaults.REDIS_PORT},
		Envs:         envs,
		VolumeMounts: map[string]string{},
	}
	deployment := factory.NewDeployment(&deploymentArgs)
	return &deployment
}

// CreateRedisService creates service for redis
func CreateRedisService(config *model.ServiceUnitConfig) *corev1.Service {
	serviceArgs := factory.ServiceArgs{
		Name:      config.Name,
		Namespace: defaults.NAMESPACE,
		Ports:     map[string]int32{"redis": defaults.REDIS_PORT},
	}
	service := factory.NewSerivce(&serviceArgs)
	return &service
}
