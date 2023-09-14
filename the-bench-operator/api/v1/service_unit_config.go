package v1

// this type is not used by the custom resource.
type ServiceUnitConfig struct {
	Name            string               `yaml:"name,omitempty" validate:"required"`
	Version         string               `yaml:"version,omitempty" validate:"required"`
	AdapterConfigs []PrimaryAdapterSpec `yaml:"adapters,omitempty" validate:"required"`
	Gateway         *Gateway             `yaml:"gateway,omitempty"`
}

// Gateway contains config information about loadgenerator
type Gateway struct {
	// VirtualUsers is the number of virtual users simulated.
	VirtualUsers int32 `yaml:"virtualUsers,omitempty"`

	// Duration given in minutes
	Duration int32 `yaml:"duration,omitempty"`
}
