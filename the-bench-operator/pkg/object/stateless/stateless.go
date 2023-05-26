package stateless

import (
	"fmt"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/object/factory"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// CreateServiceUnitDeployment creates deployment for service unit
func CreateStatelessUnitDeployment(name string) *appsv1.Deployment {
	deploymentArgs := factory.DeploymentArgs{
		Name:                   name,
		Namespace:              factory.NAMESPACE,
		Image:                  factory.SERVICE_UNIT_IMAGE,
		Replicas:               factory.REPLICAS,
		ResourceLimitsCPU:      factory.LIMIT_CPU,
		ResourceLimitsMemory:   factory.LIMIT_MEM,
		ResourceRequestsCPU:    factory.REQUEST_CPU,
		ResourceRequestsMemory: factory.REQUEST_MEM,
		Ports:                  map[string]int32{"http": factory.HTTP_PORT},
		VolumeMounts:           map[string]string{"config": "/app/config/"},
		ConfigVolume: &factory.ConfigMapVolumeArgs{
			Name: fmt.Sprintf("%s-config", name),
			Items: map[string]string{
				"config": "service-unit.yaml",
			},
		},
		// EnvVolume: &factory.ConfigMapVolumeArgs{
		// 	Name: fmt.Sprintf("%s-env", serviceUnitConfig.Name),
		// 	Items: map[string]string{
		// 		"env": ".env",
		// 	},
		// },
	}
	deployment := factory.DeploymentFactory(&deploymentArgs)
	return &deployment
}

// CreateServiceUnitService creates service for service unit
func CreateStatelessUnitService(name string) *corev1.Service {
	serviceArgs := factory.ServiceArgs{
		Name:      name,
		Namespace: factory.NAMESPACE,
		Ports:     map[string]int32{"http": factory.HTTP_PORT},
	}
	service := factory.SerivceFactory(&serviceArgs)
	return &service
}

// CreateServiceUnitConfigMap creates config config map for service unit
func CreateStatelessUnitYamlConfigMap(name string, rawConfig string) *corev1.ConfigMap {
	configMapArgs := factory.ConfigMapArgs{
		Name:      fmt.Sprintf("%s-config", name),
		Namespace: factory.NAMESPACE,
		Data: map[string]string{
			"config": rawConfig,
		},
	}
	configMap := factory.ConfigMapFactory(&configMapArgs)
	return &configMap
}

// // CreateServiceUnitEnvConfigMap creates env config map for service unit
// func CreateStatelessUnitEnvConfigMap(name string, rawEnv string) *corev1.ConfigMap {
// 	configMapArgs := factory.ConfigMapArgs{
// 		Name:      fmt.Sprintf("%s-env", name),
// 		Namespace: factory.NAMESPACE,
// 		Data: map[string]string{
// 			"env": rawEnv,
// 		},
// 	}
// 	configMap := factory.ConfigMapFactory(&configMapArgs)
// 	return &configMap
// }
