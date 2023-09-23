package factory

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ObjectMetaOptions struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
}

// NewTypeMeta create type meta for all kubernetes objects.
func NewTypeMeta(kind string, apiVersion string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       kind,
		APIVersion: apiVersion,
	}
}

// NewObjectMeta object meta data with name and namespace
func NewObjectMeta(options ObjectMetaOptions) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:        options.Name,
		Namespace:   options.Namespace,
		Labels:      options.Labels,
		Annotations: options.Annotations,
	}
}

// NewLabelSelector create lable selector
func NewLabelSelector(matchLabels map[string]string) *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: matchLabels,
	}
}

// NewLocalObjectReference create lable selector
func NewLocalObjectReference(name string) *corev1.LocalObjectReference {
	return &corev1.LocalObjectReference{Name: name}
}
