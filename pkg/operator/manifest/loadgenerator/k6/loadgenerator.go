package k6

import (
	"encoding/json"
	"fmt"

	"github.com/hanapedia/hexagon/pkg/api/defaults"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/operator/object/factory"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// CreateServiceUnitDeployment creates deployment for service unit
func CreateLoadGeneratorDeployment(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *appsv1.Deployment {
	resource := getDefaultResource()
	name := fmt.Sprintf("%s-lg", suc.Name)
	deploymentArgs := factory.DeploymentArgs{
		Name:         name,
		Namespace:    cc.Namespace,
		Annotations:  map[string]string{"rca": "ignore"},
		Image:        fmt.Sprintf("%s/%s:%s", cc.DockerHubUsername, defaults.LOAD_GENERATOR_IMAGE_NAME, suc.Version),
		Replicas:     1,
		Resource:     resource,
		Ports:        map[string]int32{"http": cc.HTTPPort},
		VolumeMounts: map[string]string{"config": "/data/"},
		ConfigVolume: &factory.ConfigMapVolumeArgs{
			Name: fmt.Sprintf("%s-config", name),
			Items: map[string]string{
				"config": "config.json",
				"routes": "routes.json",
			},
		},
		Envs:                  createLoadGeneratorDeploymentEnvs(suc),
		DisableReadinessProbe: true,
	}
	deployment := factory.NewDeployment(&deploymentArgs)

	deployment.Spec.Template.Spec.Containers[0].Command = []string{
		"k6",
		"run",
		"-o",
		"experimental-prometheus-rw",
		"/scripts/script.js",
	}

	return &deployment
}

// CreateServiceUnitService creates service for service unit
func CreateLoadGeneratorService(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *corev1.Service {
	name := fmt.Sprintf("%s-lg", suc.Name)
	serviceArgs := factory.ServiceArgs{
		Name:      name,
		Namespace: cc.Namespace,
		Ports:     map[string]int32{"http": cc.HTTPPort},
	}
	service := factory.NewSerivce(&serviceArgs)
	return &service
}

// CreateServiceUnitConfigMap creates config config map for service unit
func CreateLoadGeneratorYamlConfigMap(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *corev1.ConfigMap {
	name := fmt.Sprintf("%s-lg", suc.Name)

	rawConfig, err := generateConfigJson(suc)
	if err != nil {
		logger.Logger.Panicf("Failed to generate raw config. %s", err)
	}

	rawRoutes, err := generateRoutesJson(suc)
	if err != nil {
		logger.Logger.Panicf("Failed to generate raw route. %s", err)
	}

	configMapArgs := factory.ConfigMapArgs{
		Name:      fmt.Sprintf("%s-config", name),
		Namespace: cc.Namespace,
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
		UrlPrefix: fmt.Sprintf("http://%s:%v/", targetName, defaults.HTTP_PORT),
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

func createLoadGeneratorDeploymentEnvs(config *model.ServiceUnitConfig) []corev1.EnvVar {
	return []corev1.EnvVar{
		{Name: "TEST_NAME", Value: fmt.Sprintf("%s-lg", config.Name)},
		{Name: "K6_PROMETHEUS_RW_SERVER_URL", Value: defaults.LG_K6_PROMETHEUS_RW_SERVER_URL},
		{Name: "K6_PROMETHEUS_RW_TREND_STATS", Value: defaults.LG_K6_PROMETHEUS_RW_TREND_STATS},
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
