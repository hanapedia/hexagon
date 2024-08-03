package factory

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceArgs struct {
	Name      string
	Namespace string
	// ports are assumed to use TCP
	// ports are mapped to same target port
	Ports map[string]int32
}

// NewSerivce create new service api object
func NewSerivce(args *ServiceArgs) corev1.Service {
	return corev1.Service{
		TypeMeta: NewTypeMeta("Service", "v1"),
		ObjectMeta: NewObjectMeta(ObjectMetaOptions{
			Name:      args.Name,
			Namespace: args.Namespace,
			Labels:    map[string]string{AppLabel: args.Name},
		}),
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{AppLabel: args.Name},
			Ports:    NewServicePort(args.Ports),
		},
	}
}

// NewServicePort create service port slice
func NewServicePort(ports map[string]int32) []corev1.ServicePort {
	var containerPorts []corev1.ServicePort
	for name, port := range ports {
		containerPorts = append(containerPorts, corev1.ServicePort{
			Name:       name,
			Protocol:   corev1.ProtocolTCP,
			Port:       port,
			TargetPort: intstr.FromString(name),
		})
	}
	return containerPorts
}
