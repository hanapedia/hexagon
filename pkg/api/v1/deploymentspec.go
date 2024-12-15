package v1

import (
	corev1 "k8s.io/api/core/v1"
)

// DeploymentSpec contains configuration for deployment
type DeploymentSpec struct {
	Replicas                       int32                        `json:"replicas,omitempty"`
	Gateway                        *Gateway                     `json:"gateway,omitempty"`
	Resource                       *corev1.ResourceRequirements `json:"resources,omitempty"`
	EnvVar                         []corev1.EnvVar              `json:"env,omitempty"`
	EnableTopologySpreadConstraint bool                         `json:"enableTopologySpreadConstraint"`
	DisableReadinessProbe          bool                         `json:"disableReadinessProbe"`
}

// Gateway contains config information about loadgenerator
type Gateway struct {
	// VirtualUsers is the number of virtual users simulated.
	VirtualUsers int32 `json:"virtualUsers,omitempty"`

	// Duration given in minutes
	Duration int32 `json:"duration,omitempty"`
}
