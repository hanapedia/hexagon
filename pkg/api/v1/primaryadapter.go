package v1

import (
	"fmt"

	"github.com/hanapedia/hexagon/pkg/operator/constants"
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
	Action  constants.Action              `json:"action,omitempty" validate:"required"`
	Route   string                        `json:"route,omitempty" validate:"required"`
	Payload PayloadSpec                   `json:"payload,omitempty"`

	// applies to only gateway service
	// refers to the weight applied to the route
	// intentionally a pointer to destinguish 0
	Weight *int32 `json:"weight,omitempty"`
}

// Config fields for repository services
type RepositoryConfig struct {
	Variant constants.RepositoryVariant `json:"variant,omitempty" validate:"required,oneof=mongo redis"`
}

// Config fields for Brokers
type ConsumerConfig struct {
	Variant constants.BrokerVariant `json:"variant,omitempty" validate:"required"`
	Topic   string                  `json:"topic,omitempty" validate:"required"`
}

// Get primary adapter id
func (ias *PrimaryAdapterSpec) GetId(serviceName string) string {
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

// Get primary adapter group by key
func (ias *PrimaryAdapterSpec) GetGroupByKey() string {
	var key string
	if ias.ServerConfig != nil {
		key = ias.ServerConfig.GetGroupByKey()
	}
	if ias.ConsumerConfig != nil {
		key = ias.ConsumerConfig.GetGroupByKey()
	}
	if ias.RepositoryConfig != nil {
		key = ias.RepositoryConfig.GetGroupByKey()
	}
	return key
}

// Get invocation adapter id
func (sac *ServerConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		serviceName,
		sac.Variant,
		sac.Action,
		sac.Route,
	)
}

// Get invocation adapter id
func (sac *ServerConfig) GetGroupByKey() string {
	return string(sac.Variant)
}

// Get repository adapter id
func (sac *RepositoryConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s",
		sac.Variant,
		serviceName,
	)
}

// Get repository adapter group by key
func (rac *RepositoryConfig) GetGroupByKey() string {
	return string(rac.Variant)
}

// Get consumer adapter id
func (bac *ConsumerConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s",
		bac.Variant,
		bac.Topic,
	)
}

// Get consumer adapter id
func (bac *ConsumerConfig) GetGroupByKey() string {
	return fmt.Sprintf(
		"%s.%s",
		bac.Variant,
		bac.Topic,
	)
}

// Get consumer group id
func (bac *ConsumerConfig) GetConsumerGroupId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s.%s",
		serviceName,
		bac.Variant,
		bac.Topic,
	)
}
