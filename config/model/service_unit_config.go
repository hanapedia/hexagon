package model

type ServiceUnitConfig struct {
	Name                  string                 `yaml:"name" validate:"required"`
	IngressAdapterConfigs []IngressAdapterConfig `yaml:"ingressAdapters" validate:"required"`
}

type IngressAdapterConfig struct {
	StatelessIngressAdapterConfig *StatelessAdapterConfig `yaml:"stateless,omitempty"`
	BrokerIngressAdapterConfig    *BrokerAdapterConfig    `yaml:"broker,omitempty"`
	StatefulIngressAdapterConfig  *StatefulAdapterConfig  `yaml:"stateful,omitempty"`
	Steps                         []Step                  `yaml:"steps" validate:"required"`
}

type Step struct {
	EgressAdapterConfig *EgressAdapterConfig `yaml:"egressAdapter" validate:"required"`
	Concurrent          bool                `yaml:"concurrent,omitempty"`
}

// one of the adapter type must be provided
// Id can be used for short definition if
type EgressAdapterConfig struct {
	StatelessEgressAdapterConfig *StatelessAdapterConfig `yaml:"stateless,omitempty"`
	StatefulEgressAdapterConfig  *StatefulAdapterConfig  `yaml:"stateful,omitempty"`
	InternalEgressAdapterConfig  *InternalAdapterConfig  `yaml:"internal,omitempty"`
	BrokerEgressAdapterConfig    *BrokerAdapterConfig    `yaml:"broker,omitempty"`
	Id                           *string                 `yaml:"id,omitempty"`
}
