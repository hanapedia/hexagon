package model

import "github.com/hanapedia/the-bench/config/constants"

type ServiceUnitConfig struct {
	Name                 string                 `yaml:"name"`
	IngressAdapterConfig []IngressAdapterConfig `yaml:"ingressAdapter"`
}

type IngressAdapterConfig struct {
	StatelessIngressAdapterConfig *StatelessAdapterConfig `yaml:"stateless,omitempty"`
	BrokerIngressAdapterConfig    *BrokerAdapterConfig    `yaml:"broker,omitempty"`
	Steps    []Step                          `yaml:"steps"`
}

type Step struct {
	EgressAdapterConfig EgressAdapterConfig `yaml:"egressAdapter"`
	Concurrent          bool                `yaml:"concurrent,omitempty"`
}

// one of the adapter type must be provided
// Id can be used for short definition if
type EgressAdapterConfig struct {
	StatelessEgressAdapterConfig *StatelessAdapterConfig `yaml:"stateless,omitempty"`
	StatefulEgressAdapterConfig  *StatefulAdapterConfig  `yaml:"stateful,omitempty"`
	InternalEgressAdapterConfig  *InternalAdapterConfig  `yaml:"internal,omitempty"`
	BrokerEgressAdapterConfig    *BrokerAdapterConfig    `yaml:"broker,omitempty"`
	Id                    *string                `yaml:"id,omitempty"`
}

// Config fields for stateful services
type StatelessAdapterConfig struct {
	Variant constants.StatelessAdapterVariant `yaml:"variant"`
	Service string                           `yaml:"service,omitempty"`
	Action  constants.Action                 `yaml:"action"`
	Route   string                           `yaml:"route"`
}

// Config fields for stateful services
type StatefulAdapterConfig struct {
	Variant constants.StatefulAdapterVariant `yaml:"variant"`
	Action  constants.Action                `yaml:"action"`
	Size    string                          `yaml:"size"`
}

// Config fields for Brokers
type BrokerAdapterConfig struct {
	Variant constants.BrokerAdapterVariant `yaml:"variant"`
	Topic   string                        `yaml:"topic"`
}

// Config fields for Internal services
type InternalAdapterConfig struct {
	Resources []string `yaml:"resources"`
	Duration  string   `yaml:"duration"`
	Load      string   `yaml:"load"`
}
