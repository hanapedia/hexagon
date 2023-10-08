package k6

import (
	"encoding/json"
	"fmt"

	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
	"github.com/hanapedia/the-bench/pkg/operator/object/factory"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// CreateServiceUnitDeployment creates deployment for service unit
func CreateLoadGeneratorDeployment(config *model.ServiceUnitConfig) *appsv1.Deployment {
	resource := getDefaultResource()
	name := fmt.Sprintf("%s-lg", config.Name)
	deploymentArgs := factory.DeploymentArgs{
		Name:         name,
		Namespace:    model.NAMESPACE,
		Annotations:  map[string]string{"rca": "ignore"},
		Image:        fmt.Sprintf("%s:%s", model.LOAD_GENERATOR_IMAGE_NAME, config.Version),
		Replicas:     1,
		Resource:     resource,
		Ports:        map[string]int32{"http": model.HTTP_PORT},
		VolumeMounts: map[string]string{"config": "/data/"},
		ConfigVolume: &factory.ConfigMapVolumeArgs{
			Name: fmt.Sprintf("%s-config", name),
			Items: map[string]string{
				"config": "config.json",
				"routes": "routes.json",
			},
		},
	}
	deployment := factory.NewDeployment(&deploymentArgs)
	return &deployment
}

// CreateServiceUnitService creates service for service unit
func CreateLoadGeneratorService(config *model.ServiceUnitConfig) *corev1.Service {
	name := fmt.Sprintf("%s-lg", config.Name)
	serviceArgs := factory.ServiceArgs{
		Name:      name,
		Namespace: model.NAMESPACE,
		Ports:     map[string]int32{"http": model.HTTP_PORT},
	}
	service := factory.NewSerivce(&serviceArgs)
	return &service
}

// CreateServiceUnitConfigMap creates config config map for service unit
func CreateLoadGeneratorYamlConfigMap(config *model.ServiceUnitConfig) *corev1.ConfigMap {
	name := fmt.Sprintf("%s-lg", config.Name)

	rawConfig, err := generateConfigJson(config)
	if err != nil {
		logger.Logger.Panicf("Failed to generate raw config. %s", err)
	}

	rawRoutes, err := generateRoutesJson(config)
	if err != nil {
		logger.Logger.Panicf("Failed to generate raw route. %s", err)
	}

	configMapArgs := factory.ConfigMapArgs{
		Name:      fmt.Sprintf("%s-config", name),
		Namespace: model.NAMESPACE,
		Data: map[string]string{
			"config": string(rawConfig),
			"routes": string(rawRoutes),
		},
	}
	configMap := factory.NewConfigMap(&configMapArgs)
	return &configMap
}

// createLoadGeneratorConfig creates load generator config struct
func createLoadGeneratorConfig(vus, duration int32, targetName string) Config {
	return Config{
		Vus:       vus,
		Duration:  fmt.Sprintf("%vm", duration),
		UrlPrefix: fmt.Sprintf("http://%s:%v/", targetName, model.HTTP_PORT),
	}
}

// createLoadGeneratorRoutes creates load generator route struct
func createLoadGeneratorRoutes(route string, method constants.HttpMethod, weight int32) Route {
	return Route{
		Route:  route,
		Method: method,
		Weight: weight,
	}
}

// generateConfigJson generate json content for config
func generateConfigJson(config *model.ServiceUnitConfig) ([]byte, error) {
	if config.DeploymentSpec.Gateway == nil {
		logger.Logger.Panic("Gateway config not found.")
	}

	k6Config := createLoadGeneratorConfig(
		config.DeploymentSpec.Gateway.VirtualUsers,
		config.DeploymentSpec.Gateway.Duration,
		config.Name,
	)
	return json.Marshal(k6Config)
}

// GenerateConfigJson generate json content for config
func generateRoutesJson(config *model.ServiceUnitConfig) ([]byte, error) {
	if config.DeploymentSpec.Gateway == nil {
		logger.Logger.Panic("Gateway config not found.")
	}
	var routes []Route
	for _, primaryAdapter := range config.AdapterConfigs {
		if primaryAdapter.ServerConfig != nil {
			var weight int32 = 1
			if primaryAdapter.ServerConfig.Weight != nil {
				weight = *primaryAdapter.ServerConfig.Weight
			}
			route := createLoadGeneratorRoutes(
				primaryAdapter.ServerConfig.Route,
				constants.GetHttpMethodFromAction(primaryAdapter.ServerConfig.Action),
				weight,
			)
			routes = append(routes, route)
		}
	}
	return json.Marshal(routes)
}
