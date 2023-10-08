package k6

import (
	"github.com/hanapedia/the-bench/pkg/api/defaults"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func getDefaultResource() *corev1.ResourceRequirements {
	return &corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(defaults.LIMIT_CPU),
			corev1.ResourceMemory: resource.MustParse(defaults.LIMIT_MEM_LG),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(defaults.REQUEST_CPU),
			corev1.ResourceMemory: resource.MustParse(defaults.REQUEST_MEM),
		},
	}
}
