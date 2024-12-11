package v1

// this type is not used by the custom resource.
type ServiceUnitConfig struct {
	ConfigTemplate
	Name           string               `json:"name,omitempty" validate:"required"`
	Version        string               `json:"version,omitempty" validate:"required"`
	AdapterConfigs []*PrimaryAdapterSpec `json:"adapters,omitempty" validate:"required"`
	DeploymentSpec DeploymentSpec       `json:"deployment,omitempty"`
}
