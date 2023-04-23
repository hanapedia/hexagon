// A ingress Adapter can have multiple handlers
// A handler can have multiple tasks sets
// A task have a single egress adapter
package core

import "github.com/hanapedia/the-bench/config/constants"

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
	Register(*Handler) error
}

type IngressAdapterError struct {
	IngressAdapter *IngressAdapter
	Error          error
}

type Handler struct {
	ID       string
	Name     string
	Protocol constants.AdapterProtocol
	Action   string
	TaskSets []TaskSet
}

type TaskSet struct {
	EgressAdapter EgressAdapter
	Concurrent    bool
}
