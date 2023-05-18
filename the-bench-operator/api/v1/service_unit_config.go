package v1

// this type is not used by the custom resource.
type ServiceUnitConfig struct {
	Name                  string               `yaml:"name,omitempty" validate:"required"`
	IngressAdapterConfigs []IngressAdapterSpec `yaml:"ingressAdapters,omitempty" validate:"required"`
}
