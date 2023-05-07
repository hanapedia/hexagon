package model

type ServiceUnitConfig struct {
	Name                  string                 `yaml:"name" validate:"required"`
	IngressAdapterConfigs []IngressAdapterConfig `yaml:"ingressAdapters" validate:"required"`
}

type IngressAdapterConfig struct {
	StatelessIngressAdapterConfig *StatelessIngressAdapterConfig `yaml:"stateless,omitempty"`
	BrokerIngressAdapterConfig    *BrokerIngressAdapterConfig    `yaml:"broker,omitempty"`
	StatefulIngressAdapterConfig  *StatefulIngressAdapterConfig  `yaml:"stateful,omitempty"`
	Steps                         []Step                         `yaml:"steps" validate:"required"`
}

type Step struct {
	EgressAdapterConfig *EgressAdapterConfig `yaml:"egressAdapter" validate:"required"`
	Concurrent          bool                 `yaml:"concurrent,omitempty"`
}

// one of the adapter type must be provided
// Id can be used for short definition if
type EgressAdapterConfig struct {
	StatelessEgressAdapterConfig *StatelessEgressAdapterConfig `yaml:"stateless,omitempty"`
	StatefulEgressAdapterConfig  *StatefulEgressAdapterConfig  `yaml:"stateful,omitempty"`
	InternalEgressAdapterConfig  *InternalAdapterConfig        `yaml:"internal,omitempty"`
	BrokerEgressAdapterConfig    *BrokerEgressAdapterConfig    `yaml:"broker,omitempty"`
	Id                           *string                       `yaml:"id,omitempty"`
}
