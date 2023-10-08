package mongo

import (
	"fmt"

	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/object/factory"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// CreateServiceUnitDeployment creates deployment for service unit
func CreateMongoDeployment(config *model.ServiceUnitConfig) *appsv1.Deployment {
	replica := config.DeploymentSpec.Replicas
	if replica <= 0 {
		replica = 1
	}

	resource := config.DeploymentSpec.Resource
	if resource == nil {
		resource = getDefaultResource()
	}

	envs := config.DeploymentSpec.EnvVar
	if envs == nil {
		envs = getDefaultEnvs()
	} else {
		envs = append(envs, getDefaultEnvs()...)
	}

	deploymentArgs := factory.DeploymentArgs{
		Name:         config.Name,
		Namespace:    model.NAMESPACE,
		Annotations:  map[string]string{"rca": "ignore"},
		Image:        fmt.Sprintf("%s:%s", model.MONGO_IMAGE_NAME, config.Version),
		Replicas:     replica,
		Resource:     resource,
		Ports:        map[string]int32{"mongo": model.MONGO_PORT},
		Envs:         envs,
		VolumeMounts: map[string]string{},
	}
	deployment := factory.NewDeployment(&deploymentArgs)
	return &deployment
}

// CreateServiceUnitService creates service for service unit
func CreateMongoService(config *model.ServiceUnitConfig) *corev1.Service {
	serviceArgs := factory.ServiceArgs{
		Name:      config.Name,
		Namespace: model.NAMESPACE,
		Ports:     map[string]int32{"mongo": model.MONGO_PORT},
	}
	service := factory.NewSerivce(&serviceArgs)
	return &service
}
