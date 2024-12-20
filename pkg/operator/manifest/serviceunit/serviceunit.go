package serviceunit

import (
	"fmt"
	"slices"

	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/hexagon/pkg/api/defaults"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/object/factory"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateServiceUnitDeployment creates deployment for service unit
func CreateStatelessUnitDeployment(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *appsv1.Deployment {
	replica := suc.DeploymentSpec.Replicas
	if replica <= 0 {
		if cc.Deployment.Replicas <= 0 {
			replica = 1
		} else {
			replica = cc.Deployment.Replicas
		}
	}

	resource := suc.DeploymentSpec.Resource
	if resource == nil {
		if cc.Deployment.Resource == nil {
			resource = getDefaultResource()
		} else {
			resource = cc.Deployment.Resource
		}
	}

	deploymentArgs := factory.DeploymentArgs{
		Name:      suc.Name,
		Namespace: cc.Namespace,
		Image: fmt.Sprintf(
			"%s/%s:%s",
			cc.DockerHubUsername,
			defaults.SERVICE_UNIT_IMAGE_NAME,
			suc.Version,
		),
		Replicas:     replica,
		Resource:     resource,
		Ports:        getPorts(cc),
		VolumeMounts: map[string]string{"config": "/app/config/"},
		Envs: slices.Concat(
			config.FromClusterConfig(*cc).AsK8sEnvVars(),
			suc.DeploymentSpec.EnvVar,
		),
		ConfigVolume: &factory.ConfigMapVolumeArgs{
			Name: fmt.Sprintf("%s-config", suc.Name),
			Items: map[string]string{
				"config": constants.SERVICE_UNIT_CONFIG_FILE_NAME,
			},
		},
		EnableTopologySpreadConstraint: suc.DeploymentSpec.EnableTopologySpreadConstraint,
		DisableReadinessProbe:          suc.DeploymentSpec.DisableReadinessProbe,
		ReadinessProbePort:             intstr.FromString(HEALTH_PORT_NAME),
	}
	deployment := factory.NewDeployment(&deploymentArgs)
	return &deployment
}

// CreateServiceUnitService creates service for service unit
func CreateStatelessUnitService(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *corev1.Service {
	serviceArgs := factory.ServiceArgs{
		Name:      suc.Name,
		Namespace: cc.Namespace,
		Ports:     getPorts(cc),
	}
	service := factory.NewSerivce(&serviceArgs)
	return &service
}

// CreateServiceUnitConfigMap creates config config map for service unit
func CreateStatelessUnitYamlConfigMap(suc *model.ServiceUnitConfig, cc *model.ClusterConfig, rawConfig string) *corev1.ConfigMap {
	configMapArgs := factory.ConfigMapArgs{
		Name:      fmt.Sprintf("%s-config", suc.Name),
		Namespace: cc.Namespace,
		Data: map[string]string{
			"config": rawConfig,
		},
	}
	configMap := factory.NewConfigMap(&configMapArgs)
	return &configMap
}

func CreateServiceMonitor(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *promv1.ServiceMonitor {
	serviceMonitorArgs := factory.ServiceMonitorArgs{
		Name:      suc.Name,
		Namespace: cc.Namespace,
	}
	serviceMonitor := factory.NewServiceMonitor(&serviceMonitorArgs)
	return &serviceMonitor
}
