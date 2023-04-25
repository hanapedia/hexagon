// A ingress Adapter can have multiple handlers
// A handler can have multiple tasks sets
// A task have a single egress adapter
package core

import (
	"fmt"

	"github.com/hanapedia/the-bench/config/model"
)

// IngressAdapter provides common interface for all the ingress resources.
// Example resources include:
// - REST API server
// - gRPC server
// - Kafka consumer
//
// It is intended to represent the individual interfaces on each exteranl service,
// not the services themselves; hence the name `EgressAdapter`
type IngressAdapter interface {
	Serve() error
	Register(*IngressAdapterHandler) error
}

type IngressAdapterError struct {
	IngressAdapter *IngressAdapter
	Error          error
}

// either StatelessAdapterConfig or BrokerAdapterConfig must be defined
type IngressAdapterHandler struct {
	StatelessIngressAdapterConfig *model.StatelessAdapterConfig
	BrokerIngressAdapterConfig    *model.BrokerAdapterConfig
	TaskSets                      []TaskSet
}

type TaskSet struct {
	EgressAdapter EgressAdapter
	Concurrent    bool
}

func (iah IngressAdapterHandler) GetId() string {
	var id string
	if iah.StatelessIngressAdapterConfig != nil {
		id = fmt.Sprintf(
			"%s.%s.%s.%s",
			iah.StatelessIngressAdapterConfig.Service,
			iah.StatelessIngressAdapterConfig.Variant,
			iah.StatelessIngressAdapterConfig.Action,
			iah.StatelessIngressAdapterConfig.Route,
		)
	}
	if iah.BrokerIngressAdapterConfig != nil {
		id = fmt.Sprintf(
			"%s.%s",
			iah.BrokerIngressAdapterConfig.Variant,
			iah.BrokerIngressAdapterConfig.Topic,
		)
	}
	return id
}
