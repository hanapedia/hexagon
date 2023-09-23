package factory

import (
	corev1 "k8s.io/api/core/v1"
)

type NamespaceArgs struct {
	Name        string
	Annotations map[string]string
}

// NewNamespace create type namespace kubernetes objects.
func NewNamespace(args *NamespaceArgs) corev1.Namespace {
	return corev1.Namespace{
		TypeMeta:   NewTypeMeta("Namespace", "v1"),
		ObjectMeta: NewObjectMeta(ObjectMetaOptions{Name: args.Name, Annotations: args.Annotations}),
	}
}
