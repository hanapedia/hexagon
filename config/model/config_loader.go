package model

import "github.com/hanapedia/the-bench/config/constants"

type ConfigLoader interface {
	Load() (ServiceUnitConfig, error)
}

type ServiceUnitConfig struct {
	Name           string          `yaml:"name"`
	HandlerConfigs []HandlerConfig `yaml:"handler"`
}

type HandlerConfig struct {
	Name     string                          `yaml:"name"`
	Protocol constants.IngressAdapterVairant `yaml:"protocol"`
	Action   string                          `yaml:"action"`
	Steps    []Step                          `yaml:"flow"`
}

type Step struct {
	Adapter    Adapter `yaml:"adapter"`
	Concurrent bool    `yaml:"concurrent,omitempty"`
}

// one of the adapter type must be provided
// Id can be used for short definition if
type Adapter struct {
	StatelessEgressConfig *StatelessEgressConfig `yaml:"stateless,omitempty"`
	StatefulEgressConfig  *StatefulEgressConfig  `yaml:"stateful,omitempty"`
	InternalEgressConfig  *InternalEgressConfig  `yaml:"internal,omitempty"`
	BrokerEgressConfig    *BrokerEgressConfig    `yaml:"broker,omitempty"`
	Id                    *string                `yaml:"id,omitempty"`
}

// Config fields for stateful services
type StatelessEgressConfig struct {
	Variant constants.StatelessEgressVariant `yaml:"variant"`
	Service string                           `yaml:"service"`
	Action  constants.Action                 `yaml:"action"`
	Route   string                           `yaml:"route"`
}

// Config fields for stateful services
type StatefulEgressConfig struct {
	Variant constants.StatefulEgressVariant `yaml:"variant"`
	Action  constants.Action                `yaml:"action"`
	Size    string                          `yaml:"size"`
}

// Config fields for Brokers
type BrokerEgressConfig struct {
	Variant constants.BrokerEgressVariant `yaml:"variant"`
	Topic   string                        `yaml:"topic"`
}

// Config fields for Internal services
type InternalEgressConfig struct {
	Resources []string `yaml:"resources"`
	Duration  string   `yaml:"duration"`
	Load      string   `yaml:"load"`
}
