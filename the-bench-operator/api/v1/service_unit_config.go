package v1

// this type is not used by the custom resource.
type ServiceUnitConfig struct {
	Name                  string               `yaml:"name,omitempty" validate:"required"`
	IngressAdapterConfigs []IngressAdapterSpec `yaml:"ingressAdapters,omitempty" validate:"required"`
	Gateway               *Gateway              `yaml:"gateway,omitempty"`
}

// Gateway contains config information about loadgenerator
type Gateway struct {
	// VirtualUsers is the number of virtual users simulated.
	VirtualUsers int `yaml:"virtualUsers,omitempty"`

	// Duration given in minutes
	Duration int `yaml:"duration,omitempty"`
}
