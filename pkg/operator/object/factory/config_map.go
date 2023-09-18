package factory

import (
	corev1 "k8s.io/api/core/v1"
)

type ConfigMapArgs struct {
	Name      string
	Namespace string
	Data      map[string]string
}

// NewConfigMap create config map
func NewConfigMap(args *ConfigMapArgs) corev1.ConfigMap {
	return corev1.ConfigMap{
		TypeMeta:   NewTypeMeta("ConfigMap", "v1"),
		ObjectMeta: NewObjectMeta(ObjectMetaOptions{Name: args.Name, Namespace: args.Namespace}),
		Data:       args.Data,
	}
}
