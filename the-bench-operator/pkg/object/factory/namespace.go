package factory

import (
	corev1 "k8s.io/api/core/v1"
)

type NamespaceArgs struct {
	Name        string
	Namespace   string
	Annotations map[string]string
}

// NamespaceFactory create type namespace kubernetes objects.
func NamespaceFactory(args NamespaceArgs) corev1.Namespace {
	return corev1.Namespace{
		TypeMeta:   TypeMetaFactory("Namespace", "v1"),
		ObjectMeta: ObjectMetaFactory(ObjectMetaOptions{Name: args.Name, Namespace: args.Namespace, Annotations: args.Annotations}),
	}
}
