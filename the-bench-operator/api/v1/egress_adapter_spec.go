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

// Get secondary adapter id
func (sac SecondaryAdapterConfig) GetId() string {
	var id string
	if sac.InvocationConfig != nil {
		id = sac.InvocationConfig.GetId()
	}
	if sac.ProducerConfig != nil {
		id = sac.ProducerConfig.GetId()
	}
	if sac.RepositoryConfig != nil {
		id = sac.RepositoryConfig.GetId()
	}
	return id
}

// Get primary adapter group by key
// Get secondary adapter id
func (sac SecondaryAdapterConfig) GetGroupByKey() string {
	var key string
	if sac.InvocationConfig != nil {
		key = sac.InvocationConfig.GetGroupByKey()
	}
	if sac.ProducerConfig != nil {
		key = sac.ProducerConfig.GetGroupByKey()
	}
	if sac.RepositoryConfig != nil {
		key = sac.RepositoryConfig.GetGroupByKey()
	}
	return key
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

// Get group key
func (sac InvocationConfig) GetGroupByKey() string {
	return fmt.Sprintf(
		"%s",
		sac.Variant,
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

// Get repository secondary adapter group by key
func (sac RepositoryClientConfig) GetGroupByKey() string {
	return sac.GetId()
}

// Get broker secondary adapter id
func (bac ProducerConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s",
		bac.Variant,
		bac.Topic,
	)
}

// Get broker secondary adapter group by key
func (bac ProducerConfig) GetGroupByKey() string {
	return bac.GetId()
}
