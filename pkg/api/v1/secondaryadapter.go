package v1

import (
	"fmt"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
)

// secondary Adapter definition for a step.
// one of the adapter type must be provided
type SecondaryAdapterConfig struct {
	InvocationConfig *InvocationConfig       `json:"invocation,omitempty"`
	RepositoryConfig *RepositoryClientConfig `json:"repository,omitempty"`
	ProducerConfig   *ProducerConfig         `json:"producer,omitempty"`
	StressorConfig   *StressorConfig         `json:"stressor,omitempty"`
	Id               *string                 `json:"id,omitempty"`
}

// Config fields for server services
type InvocationConfig struct {
	Variant constants.SeverAdapterVariant `json:"variant,omitempty" validate:"required,oneof=rest grpc"`
	Service string                        `json:"service,omitempty" validate:"required"`
	Action  constants.Action              `json:"action,omitempty" validate:"required,oneof=read write"`
	Route   string                        `json:"route,omitempty" validate:"required"`
	Payload constants.PayloadSizeVariant  `json:"payload,omitempty" validate:"omitempty,oneof=small medium large"`

	// applies to only clientStreams and biStream via grpc
	PayloadCount int `json:"payloadCount,omitempty"`
}

// Config fields for repository services
type RepositoryClientConfig struct {
	Name    string                       `json:"name,omitempty" validate:"required"`
	Variant constants.RepositoryVariant  `json:"variant,omitempty" validate:"required,oneof=mongo redis postgre"`
	Action  constants.Action             `json:"action,omitempty" validate:"omitempty,oneof=read write"`
	Payload constants.PayloadSizeVariant `json:"payload,omitempty" validate:"omitempty,oneof=small medium large"`
}

// Config fields for Brokers
type ProducerConfig struct {
	Variant constants.BrokerVariant      `json:"variant,omitempty" validate:"required,oneof=kafka rabbitmq pulsar"`
	Topic   string                       `json:"topic,omitempty" validate:"required"`
	Payload constants.PayloadSizeVariant `json:"payload,omitempty" validate:"omitempty,oneof=small medium large"`
}

// Config fields for Stressor
type StressorConfig struct {
	Name        string                       `json:"name,omitempty" validate:"required"`
	Variant     constants.StressorValiant    `json:"variant,omitempty" validate:"required,oneof=cpu memory disk"`
	Duration    string                       `json:"duration,omitempty" validate:"required,oneof=small medium large"`
	ThreadCount int                          `json:"threads,omitempty" validate:"omitempty"`
	Payload     constants.PayloadSizeVariant `json:"payload,omitempty" validate:"omitempty,oneof=small medium large"`
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

// Get internal adapter id
func (iac StressorConfig) GetId() string {
	return iac.Name
}
