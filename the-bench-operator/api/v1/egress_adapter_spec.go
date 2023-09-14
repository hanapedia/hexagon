package v1

import (
	"fmt"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
)

// secondary Adapter definition for a step.
// one of the adapter type must be provided
type SecondaryAdapterConfig struct {
	InvocationConfig *InvocationConfig       `json:"invocation,omitempty"`
	RepositoryConfig *RepositoryClientConfig `json:"repository,omitempty"`
	ProducerConfig   *ProducerConfig         `json:"producer,omitempty"`
	Id               *string                 `json:"id,omitempty"`
	// InternalEgressAdapterConfig  *InternalAdapterConfig        `json:"internal,omitempty"`
}

// Get secondary adapter id
func (eac SecondaryAdapterConfig) GetId() string {
	var id string
	if eac.InvocationConfig != nil {
		id = eac.InvocationConfig.GetId()
	}
	if eac.ProducerConfig != nil {
		id = eac.ProducerConfig.GetId()
	}
	if eac.RepositoryConfig != nil {
		id = eac.RepositoryConfig.GetId()
	}
	return id
}

// Config fields for server services
type InvocationConfig struct {
	Variant constants.SeverAdapterVariant `json:"variant,omitempty" validate:"required,oneof=rest grpc"`
	Service string                        `json:"service,omitempty" validate:"required"`
	Action  constants.Action              `json:"action,omitempty" validate:"required,oneof=read write"`
	Route   string                        `json:"route,omitempty" validate:"required"`
}

// Config fields for repository services
type RepositoryClientConfig struct {
	Name    string                      `json:"name,omitempty" validate:"required"`
	Variant constants.RepositoryVariant `json:"variant,omitempty" validate:"required,oneof=mongo postgre"`
	Action  constants.Action            `json:"action,omitempty" validate:"omitempty,oneof=read write"`
	Size    string                      `json:"size,omitempty" validate:"omitempty,oneof=small medium large"`
}

// Config fields for Brokers
type ProducerConfig struct {
	Variant constants.BrokerVariant `json:"variant,omitempty" validate:"required,oneof=kafka rabbitmq pulsar"`
	Topic   string                  `json:"topic,omitempty" validate:"required"`
}

// Get server secondary adapter id
func (sac InvocationConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		sac.Service,
		sac.Variant,
		sac.Action,
		sac.Route,
	)
}

// Get repository secondary adapter id
func (sac RepositoryClientConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s",
		sac.Variant,
		sac.Name,
	)
}

// Get broker secondary adapter id
func (bac ProducerConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s",
		bac.Variant,
		bac.Topic,
	)
}
