package factory

import (
	corev1 "k8s.io/api/core/v1"
)

type ConfigMapArgs struct {
	Name      string
	Namespace string
	Data      map[string]string
}

// ConfigMapFactory create config map
func ConfigMapFactory(args *ConfigMapArgs) corev1.ConfigMap {
	return corev1.ConfigMap{
		TypeMeta:   TypeMetaFactory("ConfigMap", "v1"),
		ObjectMeta: ObjectMetaFactory(ObjectMetaOptions{Name: args.Name, Namespace: args.Namespace}),
		Data:       args.Data,
	}
}
