package serviceunit

import (
	"github.com/hanapedia/hexagon/pkg/api/defaults"
	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/object/factory"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	HTTP_PORT_NAME   = "http"
	GRPC_PORT_NAME   = "grpc"
	HEALTH_PORT_NAME = "health"
)

func getDefaultResource() *corev1.ResourceRequirements {
	return &corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(defaults.LIMIT_CPU),
			corev1.ResourceMemory: resource.MustParse(defaults.LIMIT_MEM),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(defaults.REQUEST_CPU),
			corev1.ResourceMemory: resource.MustParse(defaults.REQUEST_MEM),
		},
	}
}

func getPorts(clusterConfig *v1.ClusterConfig) map[string]int32 {
	var health int32 = defaults.HEALTH_PORT
	var http int32 = defaults.HTTP_PORT
	var grpc int32 = defaults.GRPC_PORT
	var metrics int32 = defaults.METRICS_PORT

	if clusterConfig.HealthPort != 0 {
		health = clusterConfig.HealthPort
	}
	if clusterConfig.HTTPPort != 0 {
		http = clusterConfig.HTTPPort
	}
	if clusterConfig.GRPCPort != 0 {
		grpc = clusterConfig.GRPCPort
	}
	if clusterConfig.MetricsPort != 0 {
		metrics = clusterConfig.MetricsPort
	}

	return map[string]int32{
		HEALTH_PORT_NAME:          health,
		HTTP_PORT_NAME:            http,
		GRPC_PORT_NAME:            grpc,
		factory.METRICS_PORT_NAME: metrics,
	}
}
