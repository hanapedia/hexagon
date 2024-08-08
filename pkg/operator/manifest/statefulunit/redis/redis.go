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
func CreateRedisDeployment(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *appsv1.Deployment {
	replica := suc.DeploymentSpec.Replicas
	if replica <= 0 {
		replica = 1
	}

	resource := suc.DeploymentSpec.Resource
	if resource == nil {
		resource = getDefaultResource()
	}

	envs := suc.DeploymentSpec.EnvVar

	deploymentArgs := factory.DeploymentArgs{
		Name:         suc.Name,
		Namespace:    cc.Namespace,
		Annotations:  map[string]string{"rca": "ignore"},
		Image:        fmt.Sprintf("%s/%s:%s", cc.Redis.ImageName, defaults.REDIS_IMAGE_NAME, suc.Version),
		Replicas:     replica,
		Resource:     resource,
		Ports:        map[string]int32{"redis": cc.Redis.Port},
		Envs:         envs,
		VolumeMounts: map[string]string{},
	}
	deployment := factory.NewDeployment(&deploymentArgs)
	return &deployment
}

// CreateRedisService creates service for redis
func CreateRedisService(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *corev1.Service {
	serviceArgs := factory.ServiceArgs{
		Name:      suc.Name,
		Namespace: cc.Namespace,
		Ports:     map[string]int32{"redis": cc.Redis.Port},
	}
	service := factory.NewSerivce(&serviceArgs)
	return &service
}
