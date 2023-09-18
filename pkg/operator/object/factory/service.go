package factory

import (
	corev1 "k8s.io/api/core/v1"
)

type ServiceArgs struct {
	Name      string
	Namespace string
	// ports are assumed to use TCP
	// ports are mapped to same target port
	Ports map[string]int32
}

func SerivceFactory(args *ServiceArgs) corev1.Service {
	return corev1.Service{
		TypeMeta:   TypeMetaFactory("Service", "v1"),
		ObjectMeta: ObjectMetaFactory(ObjectMetaOptions{Name: args.Name, Namespace: args.Namespace}),
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": args.Name},
			Ports:    ServicePortFactory(args.Ports),
		},
	}
}

// ServicePortFactory create service port slice
func ServicePortFactory(ports map[string]int32) []corev1.ServicePort {
	var containerPorts []corev1.ServicePort
	for name, port := range ports {
		containerPorts = append(containerPorts, corev1.ServicePort{
			Name:     name,
			Protocol: corev1.ProtocolTCP,
			Port:     port,
		})
	}
	return containerPorts
}
