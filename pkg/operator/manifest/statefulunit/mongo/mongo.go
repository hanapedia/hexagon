package mongo

import (
	"fmt"

	"github.com/hanapedia/hexagon/pkg/api/defaults"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/object/factory"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// CreateServiceUnitDeployment creates deployment for service unit
func CreateMongoDeployment(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *appsv1.Deployment {
	replica := suc.DeploymentSpec.Replicas
	if replica <= 0 {
		replica = 1
	}

	resource := suc.DeploymentSpec.Resource
	if resource == nil {
		resource = getDefaultResource()
	}

	envs := suc.DeploymentSpec.EnvVar
	if envs == nil {
		envs = getDefaultEnvs()
	} else {
		envs = append(envs, getDefaultEnvs()...)
	}

	deploymentArgs := factory.DeploymentArgs{
		Name:         suc.Name,
		Namespace:    cc.Namespace,
		Annotations:  map[string]string{"rca": "ignore"},
		Image:        fmt.Sprintf("%s/%s:%s", cc.DockerHubUsername, defaults.MONGO_IMAGE_NAME, suc.Version),
		Replicas:     replica,
		Resource:     resource,
		Ports:        map[string]int32{"mongo": cc.Mongo.Port},
		Envs:         envs,
		VolumeMounts: map[string]string{},
	}
	deployment := factory.NewDeployment(&deploymentArgs)
	return &deployment
}

// CreateServiceUnitService creates service for service unit
func CreateMongoService(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *corev1.Service {
	serviceArgs := factory.ServiceArgs{
		Name:      suc.Name,
		Namespace: cc.Namespace,
		Ports:     map[string]int32{"mongo": cc.Mongo.Port},
	}
	service := factory.NewSerivce(&serviceArgs)
	return &service
}
