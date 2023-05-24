package factory

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TypeMetaFactory create type meta for all kubernetes objects.
func TypeMetaFactory(kind string, apiVersion string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       kind,
		APIVersion: apiVersion,
	}
}

// ObjectMetaFactory object meta data with name and namespace
func ObjectMetaFactory(name string, namespace string, label map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    label,
	}
}

// LabelSelectorFactory create lable selector
func LabelSelectorFactory(matchLabels map[string]string) *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: matchLabels,
	}
}

// LocalObjectReferenceFactory create lable selector
func LocalObjectReferenceFactory(name string) *corev1.LocalObjectReference {
	return &corev1.LocalObjectReference{Name: name}
}
