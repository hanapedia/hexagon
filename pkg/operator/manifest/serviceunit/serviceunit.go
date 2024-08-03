package serviceunit

import (
	"fmt"

	"github.com/hanapedia/hexagon/pkg/api/defaults"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/object/factory"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// CreateServiceUnitDeployment creates deployment for service unit
func CreateStatelessUnitDeployment(config *model.ServiceUnitConfig) *appsv1.Deployment {
	replica := config.DeploymentSpec.Replicas
	if replica <= 0 {
		replica = 1
	}

	resource := config.DeploymentSpec.Resource
	if resource == nil {
		resource = getDefaultResource()
	}

	deploymentArgs := factory.DeploymentArgs{
		Name:         config.Name,
		Namespace:    defaults.NAMESPACE,
		Image:        fmt.Sprintf("%s:%s", defaults.SERVICE_UNIT_IMAGE_NAME, config.Version),
		Replicas:     replica,
		Resource:     resource,
		Ports:        getDefaultPorts(),
		VolumeMounts: map[string]string{"config": "/app/config/"},
		Envs:         config.DeploymentSpec.EnvVar,
		ConfigVolume: &factory.ConfigMapVolumeArgs{
			Name: fmt.Sprintf("%s-config", config.Name),
			Items: map[string]string{
				"config": "service-unit.yaml",
			},
		},
	}
	deployment := factory.NewDeployment(&deploymentArgs)
	return &deployment
}

// CreateServiceUnitService creates service for service unit
func CreateStatelessUnitService(config *model.ServiceUnitConfig) *corev1.Service {
	serviceArgs := factory.ServiceArgs{
		Name:      config.Name,
		Namespace: defaults.NAMESPACE,
		Ports:     getDefaultPorts(),
	}
	service := factory.NewSerivce(&serviceArgs)
	return &service
}

// CreateServiceUnitConfigMap creates config config map for service unit
func CreateStatelessUnitYamlConfigMap(config *model.ServiceUnitConfig, rawConfig string) *corev1.ConfigMap {
	configMapArgs := factory.ConfigMapArgs{
		Name:      fmt.Sprintf("%s-config", config.Name),
		Namespace: defaults.NAMESPACE,
		Data: map[string]string{
			"config": rawConfig,
		},
	}
	configMap := factory.NewConfigMap(&configMapArgs)
	return &configMap
}

func CreateServiceMonitor(config *model.ServiceUnitConfig) *promv1.ServiceMonitor {
	serviceMonitorArgs := factory.ServiceMonitorArgs{
		Name:      config.Name,
		Namespace: defaults.NAMESPACE,
	}
	serviceMonitor := factory.NewServiceMonitor(&serviceMonitorArgs)
	return &serviceMonitor
}
