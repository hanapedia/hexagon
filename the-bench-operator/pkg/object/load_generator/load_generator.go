package loadgenerator

import (
	"fmt"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/object/factory"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// CreateServiceUnitDeployment creates deployment for service unit
func CreateLoadGeneratorDeployment(name string) *appsv1.Deployment {
	deploymentArgs := factory.DeploymentArgs{
		Name:                   name,
		Namespace:              factory.NAMESPACE,
		Image:                  factory.LOAD_GENERATOR_IMAGE,
		Replicas:               factory.REPLICAS,
		ResourceLimitsCPU:      factory.LIMIT_CPU,
		ResourceLimitsMemory:   factory.LIMIT_MEM_LG,
		ResourceRequestsCPU:    factory.REQUEST_CPU,
		ResourceRequestsMemory: factory.REQUEST_MEM,
		Ports:                  map[string]int32{"http": factory.HTTP_PORT},
		VolumeMounts:           map[string]string{"config": "/data/"},
		ConfigVolume: &factory.ConfigMapVolumeArgs{
			Name: fmt.Sprintf("%s-config", name),
			Items: map[string]string{
				"config": "config.json",
				"routes": "routes.json",
			},
		},
	}
	deployment := factory.DeploymentFactory(&deploymentArgs)
	return &deployment
}

// CreateServiceUnitService creates service for service unit
func CreateLoadGeneratorService(name string) *corev1.Service {
	serviceArgs := factory.ServiceArgs{
		Name:      name,
		Namespace: factory.NAMESPACE,
		Ports:     map[string]int32{"http": factory.HTTP_PORT},
	}
	service := factory.SerivceFactory(&serviceArgs)
	return &service
}

// CreateServiceUnitConfigMap creates config config map for service unit
func CreateLoadGeneratorYamlConfigMap(name, rawConfig, rawRoutes string) *corev1.ConfigMap {
	configMapArgs := factory.ConfigMapArgs{
		Name:      fmt.Sprintf("%s-config", name),
		Namespace: factory.NAMESPACE,
		Data: map[string]string{
			"config": rawConfig,
			"routes": rawRoutes,
		},
	}
	configMap := factory.ConfigMapFactory(&configMapArgs)
	return &configMap
}

// CreateLoadGeneratorConfig creates load generator config struct
func CreateLoadGeneratorConfig(vus, duration int32, targetName string) Config {
	return Config{
		Vus:       vus,
		Duration:  fmt.Sprintf("%vm", duration),
		UrlPrefix: fmt.Sprintf("http://%s:%v/", targetName, factory.HTTP_PORT),
	}
}

// CreateLoadGeneratorRoutes creates load generator route struct
func CreateLoadGeneratorRoutes(route string, method constants.HttpMethod, weight int32) Route {
	return Route{
		Route:  route,
		Method: method,
		Weight: weight,
	}
}
