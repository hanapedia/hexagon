package v1

// this type is not used by the custom resource.
type ServiceUnitConfig struct {
	Name           string               `json:"name,omitempty" validate:"required"`
	Version        string               `json:"version,omitempty" validate:"required"`
	AdapterConfigs []PrimaryAdapterSpec `json:"adapters,omitempty" validate:"required"`
	Gateway        *Gateway             `json:"gateway,omitempty"`
}

// Gateway contains config information about loadgenerator
type Gateway struct {
	// VirtualUsers is the number of virtual users simulated.
	VirtualUsers int32 `json:"virtualUsers,omitempty"`

	// Duration given in minutes
	Duration int32 `json:"duration,omitempty"`
}
