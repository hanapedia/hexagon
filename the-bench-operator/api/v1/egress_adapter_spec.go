package v1

import (
	"fmt"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
)

// Config fields for stateful services
type StatelessEgressAdapterConfig struct {
	Variant constants.StatelessAdapterVariant `json:"variant,omitempty" yaml:"variant,omitempty" validate:"required,oneof=rest grpc"`
	Service string                            `json:"service,omitempty" yaml:"service,omitempty" validate:"required"`
	Action  constants.Action                  `json:"action,omitempty" yaml:"action,omitempty" validate:"required,oneof=read write"`
	Route   string                            `json:"route,omitempty" yaml:"route,omitempty" validate:"required"`
}

func (sac StatelessEgressAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		sac.Service,
		sac.Variant,
		sac.Action,
		sac.Route,
	)
}

// Config fields for stateful services
type StatefulEgressAdapterConfig struct {
	Name    string                           `json:"name,omitempty" yaml:"name,omitempty" validate:"required"`
	Variant constants.StatefulAdapterVariant `json:"variant,omitempty" yaml:"variant,omitempty" validate:"required,oneof=mongo postgre"`
	Action  constants.Action                 `json:"action,omitempty" yaml:"action,omitempty" validate:"omitempty,oneof=read write"`
	Size    string                           `json:"size,omitempty" yaml:"size,omitempty" validate:"omitempty,oneof=small medium large"`
}

func (sac StatefulEgressAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s",
		sac.Variant,
		sac.Name,
	)
}

// Config fields for Brokers
type BrokerEgressAdapterConfig struct {
	Variant constants.BrokerAdapterVariant `json:"variant,omitempty" yaml:"variant,omitempty" validate:"required,oneof=kafka rabbitmq pulsar"`
	Topic   string               `json:"topic,omitempty" yaml:"topic,omitempty" validate:"required"`
}

func (bac BrokerEgressAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s",
		bac.Variant,
		bac.Topic,
	)
}
