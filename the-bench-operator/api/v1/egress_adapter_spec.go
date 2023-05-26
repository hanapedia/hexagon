package v1

import (
	"fmt"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
)

// egress adapter definition for a step.
// one of the adapter type must be provided
type EgressAdapterConfig struct {
	StatelessEgressAdapterConfig *StatelessEgressAdapterConfig `json:"stateless,omitempty" yaml:"stateless,omitempty"`
	StatefulEgressAdapterConfig  *StatefulEgressAdapterConfig  `json:"stateful,omitempty" yaml:"stateful,omitempty"`
	InternalEgressAdapterConfig  *InternalAdapterConfig        `json:"internal,omitempty" yaml:"internal,omitempty"`
	BrokerEgressAdapterConfig    *BrokerEgressAdapterConfig    `json:"broker,omitempty" yaml:"broker,omitempty"`
	Id                           *string                       `json:"id,omitempty" yaml:"id,omitempty"`
}

// Get egress adapter id
func (eac EgressAdapterConfig) GetId() string {
	var id string
	if eac.StatelessEgressAdapterConfig != nil {
		id = eac.StatelessEgressAdapterConfig.GetId()
	}
	if eac.BrokerEgressAdapterConfig != nil {
		id = eac.BrokerEgressAdapterConfig.GetId()
	}
	if eac.StatefulEgressAdapterConfig != nil {
		id = eac.StatefulEgressAdapterConfig.GetId()
	}
	return id
}


// Config fields for stateful services
type StatelessEgressAdapterConfig struct {
	Variant constants.StatelessAdapterVariant `json:"variant,omitempty" yaml:"variant,omitempty" validate:"required,oneof=rest grpc"`
	Service string                            `json:"service,omitempty" yaml:"service,omitempty" validate:"required"`
	Action  constants.Action                  `json:"action,omitempty" yaml:"action,omitempty" validate:"required,oneof=read write"`
	Route   string                            `json:"route,omitempty" yaml:"route,omitempty" validate:"required"`
}

// Config fields for stateful services
type StatefulEgressAdapterConfig struct {
	Name    string                           `json:"name,omitempty" yaml:"name,omitempty" validate:"required"`
	Variant constants.StatefulAdapterVariant `json:"variant,omitempty" yaml:"variant,omitempty" validate:"required,oneof=mongo postgre"`
	Action  constants.Action                 `json:"action,omitempty" yaml:"action,omitempty" validate:"omitempty,oneof=read write"`
	Size    string                           `json:"size,omitempty" yaml:"size,omitempty" validate:"omitempty,oneof=small medium large"`
}

// Config fields for Brokers
type BrokerEgressAdapterConfig struct {
	Variant constants.BrokerAdapterVariant `json:"variant,omitempty" yaml:"variant,omitempty" validate:"required,oneof=kafka rabbitmq pulsar"`
	Topic   string                         `json:"topic,omitempty" yaml:"topic,omitempty" validate:"required"`
}

// Get stateless egress adapter id
func (sac StatelessEgressAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		sac.Service,
		sac.Variant,
		sac.Action,
		sac.Route,
	)
}

// Get stateful egress adapter id
func (sac StatefulEgressAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s",
		sac.Variant,
		sac.Name,
	)
}

// Get broker egress adapter id
func (bac BrokerEgressAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s",
		bac.Variant,
		bac.Topic,
	)
}
