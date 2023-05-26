package v1

import (
	"fmt"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ingress adapter
// must be attachted to a service unit
type IngressAdapterSpec struct {
	StatelessIngressAdapterConfig *StatelessIngressAdapterConfig `json:"stateless,omitempty" yaml:"stateless,omitempty"`
	BrokerIngressAdapterConfig    *BrokerIngressAdapterConfig    `json:"broker,omitempty" yaml:"broker,omitempty"`
	StatefulIngressAdapterConfig  *StatefulIngressAdapterConfig  `json:"stateful,omitempty" yaml:"stateful,omitempty"`
	Steps                         []Step                         `json:"steps,omitempty" yaml:"steps,omitempty" validate:"required"`

	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

// A task to be performed in a single step
type Step struct {
	EgressAdapterConfig *EgressAdapterConfig `json:"egressAdapter,omitempty" yaml:"egressAdapter,omitempty" validate:"required"`
	Concurrent          bool                 `json:"concurrent,omitempty" yaml:"concurrent,omitempty"`
}

// Config fields for stateful services
type StatelessIngressAdapterConfig struct {
	Variant constants.StatelessAdapterVariant `json:"variant,omitempty" yaml:"variant,omitempty" validate:"required,oneof=rest grpc"`
	Action  constants.Action                  `json:"action,omitempty" yaml:"action,omitempty" validate:"required,oneof=read write"`
	Route   string                            `json:"route,omitempty" yaml:"route,omitempty" validate:"required"`
	// applies to only gateway service
	// refers to the weight applied to the route
	// intentionally a pointer to destinguish 0
	Weight  *int32                               `json:"weight,omitempty" yaml:"weight,omitempty"`
}

// Config fields for stateful services
type StatefulIngressAdapterConfig struct {
	Variant constants.StatefulAdapterVariant `json:"variant,omitempty" yaml:"variant,omitempty" validate:"required,oneof=mongo postgre"`
}

// Config fields for Brokers
type BrokerIngressAdapterConfig struct {
	Variant constants.BrokerAdapterVariant `json:"variant,omitempty" yaml:"variant,omitempty" validate:"required,oneof=kafka rabbitmq pulsar"`
	Topic   string                         `json:"topic,omitempty" yaml:"topic,omitempty" validate:"required"`
}

// Config fields for Internal services
type InternalAdapterConfig struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty" validate:"required"`
	Resource string `json:"resource,omitempty" yaml:"resource,omitempty" validate:"required,oneof=cpu memory disk network"`
	Duration string `json:"duration,omitempty" yaml:"duration,omitempty" validate:"required,oneof=small medium large"`
	Load     string `json:"load,omitempty" yaml:"load,omitempty" validate:"required,oneof=small medium large"`
}

// Get ingress adapter id
func (ias IngressAdapterSpec) GetId(serviceName string) string {
	var id string
	if ias.StatelessIngressAdapterConfig != nil {
		id = ias.StatelessIngressAdapterConfig.GetId(serviceName)
	}
	if ias.BrokerIngressAdapterConfig != nil {
		id = ias.BrokerIngressAdapterConfig.GetId(serviceName)
	}
	if ias.StatefulIngressAdapterConfig != nil {
		id = ias.StatefulIngressAdapterConfig.GetId(serviceName)
	}
	return id
}

// Get stateless ingress adapter id
func (sac StatelessIngressAdapterConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		serviceName,
		sac.Variant,
		sac.Action,
		sac.Route,
	)
}

// Get stateful ingress adapter id
func (sac StatefulIngressAdapterConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s",
		sac.Variant,
		serviceName,
	)
}

// Get broker ingress adapter id
func (bac BrokerIngressAdapterConfig) GetId(serviceName string) string {
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
