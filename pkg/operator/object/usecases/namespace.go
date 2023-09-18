package usecases

import (
	"github.com/hanapedia/the-bench/pkg/operator/object/factory"

	corev1 "k8s.io/api/core/v1"
)

// CreateNamespace creates namespace for service unit
func CreateNamespace() *corev1.Namespace {
	namespaceArgs := factory.NamespaceArgs{
		Name:        factory.NAMESPACE,
		Annotations: map[string]string{"linkerd.io/inject": "enabled"},
	}
	namespace := factory.NamespaceFactory(&namespaceArgs)
	return &namespace
}
