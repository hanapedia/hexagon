package mongo

import (
	"github.com/hanapedia/the-bench/pkg/api/defaults"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

func getDefaultEnvs() []corev1.EnvVar {
	return []corev1.EnvVar{
		// {
		// 	Name: "MONGO_INITDB_ROOT_USERNAME",
		// 	Value: "root",
		// },
		// {
		// 	Name: "MONGO_INITDB_ROOT_PASSWORD",
		// 	Value: "password",
		// },
		{
			Name: "MONGO_INITDB_DATABASE",
			Value: "mongo",
		},
	}
}
