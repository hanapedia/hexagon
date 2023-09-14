package v1

import (
	"fmt"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PrimaryAdapterSpec
// must be attachted to a service unit
type PrimaryAdapterSpec struct {
	ServerConfig     *ServerConfig     `json:"server,omitempty"`
	ConsumerConfig   *ConsumerConfig   `json:"consumer,omitempty"`
	RepositoryConfig *RepositoryConfig `json:"repository,omitempty"`
	Steps            []Step            `json:"steps,omitempty" validate:"required"`

	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

// A task to be performed in a single step
type Step struct {
	AdapterConfig *SecondaryAdapterConfig `json:"adapter,omitempty" validate:"required"`
	Concurrent    bool                    `json:"concurrent,omitempty" `
}

// Config fields for repository services
type ServerConfig struct {
	Variant constants.SeverAdapterVariant `json:"variant,omitempty" validate:"required,oneof=rest grpc"`
	Action  constants.Action                  `json:"action,omitempty" validate:"required,oneof=read write"`
	Route   string                            `json:"route,omitempty" validate:"required"`
	// applies to only gateway service
	// refers to the weight applied to the route
	// intentionally a pointer to destinguish 0
	Weight *int32 `json:"weight,omitempty"`
}

// Config fields for repository services
type RepositoryConfig struct {
	Variant constants.RepositoryVariant `json:"variant,omitempty" validate:"required,oneof=mongo postgre"`
}

// Config fields for Brokers
type ConsumerConfig struct {
	Variant constants.BrokerVariant `json:"variant,omitempty" validate:"required,oneof=kafka rabbitmq pulsar"`
	Topic   string                         `json:"topic,omitempty" validate:"required"`
}

// Config fields for Internal services
type InternalAdapterConfig struct {
	Name     string `json:"name,omitempty" validate:"required"`
	Resource string `json:"resource,omitempty" validate:"required,oneof=cpu memory disk network"`
	Duration string `json:"duration,omitempty" validate:"required,oneof=small medium large"`
	Load     string `json:"load,omitempty" validate:"required,oneof=small medium large"`
}

// Get primary adapter id
func (ias PrimaryAdapterSpec) GetId(serviceName string) string {
	var id string
	if ias.ServerConfig != nil {
		id = ias.ServerConfig.GetId(serviceName)
	}
	if ias.ConsumerConfig != nil {
		id = ias.ConsumerConfig.GetId(serviceName)
	}
	if ias.RepositoryConfig != nil {
		id = ias.RepositoryConfig.GetId(serviceName)
	}
	return id
}

// Get invocation adapter id
func (sac ServerConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		serviceName,
		sac.Variant,
		sac.Action,
		sac.Route,
	)
}

// Get repository adapter id
func (sac RepositoryConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s",
		sac.Variant,
		serviceName,
	)
}

// Get consumer adapter id
func (bac ConsumerConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s",
		bac.Variant,
		bac.Topic,
	)
}

// Get internal adapter id
func (iac InternalAdapterConfig) GetId() string {
	return iac.Name
}
