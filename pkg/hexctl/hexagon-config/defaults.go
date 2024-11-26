package hexagonconfig

import (
	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	DEFAULT_GET_ROUTE = "get"
	DEFAULT_GW_ROUTE  = "gateway"
)

/* deployment: */
/*   replicas: 1 */
/*   resources: */
/*     limits: */
/*       cpu: 500m */
/*       memory: 1Gi */
/*     requests: */
/*       cpu: 500m */
/*       memory: 1Gi */
func NewDefaultDeploymentSpec() v1.DeploymentSpec {
	return v1.DeploymentSpec{
		Replicas: 1,
		Resource: &corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
		},
	}
}
