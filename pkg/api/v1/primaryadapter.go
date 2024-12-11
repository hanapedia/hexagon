package v1

import (
	"fmt"

	"github.com/hanapedia/hexagon/pkg/operator/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PrimaryAdapterType int64

const (
	Server PrimaryAdapterType = iota
	Consumer
	Repository
)

// PrimaryAdapterSpec
// must be attachted to a service unit
type PrimaryAdapterSpec struct {
	ServerConfig     *ServerConfig     `json:"server,omitempty"`
	ConsumerConfig   *ConsumerConfig   `json:"consumer,omitempty"`
	RepositoryConfig *RepositoryConfig `json:"repository,omitempty"`
	TaskSpecs        []*TaskSpec        `json:"tasks,omitempty" validate:"required"`

	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

// A spec for atask to be performed
type TaskSpec struct {
	AdapterConfig *SecondaryAdapterConfig `json:"adapter,omitempty" validate:"required"`
	Concurrent    bool                    `json:"concurrent,omitempty" `
	Resiliency    ResiliencySpec          `json:"resiliency,omitempty"`
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
func (pas *PrimaryAdapterSpec) GetId(serviceName string) string {
	var id string
	switch pas.Type() {
	case Server:
		id = pas.ServerConfig.GetId(serviceName)
	case Consumer:
		id = pas.ConsumerConfig.GetId(serviceName)
	case Repository:
		id = pas.RepositoryConfig.GetId(serviceName)
	}
	return id
}

// Get primary adapter group by key
func (pas *PrimaryAdapterSpec) GetGroupByKey() string {
	var key string
	switch pas.Type() {
	case Server:
		key = pas.ServerConfig.GetGroupByKey()
	case Consumer:
		key = pas.ConsumerConfig.GetGroupByKey()
	case Repository:
		key = pas.RepositoryConfig.GetGroupByKey()
	}
	return key
}

// Get primary adapter type
func (pas *PrimaryAdapterSpec) Type() PrimaryAdapterType {
	if pas.ServerConfig != nil {
		return Server
	}
	if pas.ConsumerConfig != nil {
		return Consumer
	}
	if pas.RepositoryConfig != nil {
		return Repository
	}

	// Should not happen
	return 0
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
