package serviceunit

import (
	"github.com/hanapedia/hexagon/pkg/api/defaults"
	"github.com/hanapedia/hexagon/pkg/operator/object/factory"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	HTTP_PORT_NAME = "http"
	GRPC_PORT_NAME = "grpc"
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

func getDefaultPorts() map[string]int32 {
	return map[string]int32{
		HTTP_PORT_NAME:            defaults.HTTP_PORT,
		GRPC_PORT_NAME:            defaults.GRPC_PORT,
		factory.METRICS_PORT_NAME: defaults.METRICS_PORT,
	}
}
