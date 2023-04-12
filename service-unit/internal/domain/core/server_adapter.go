// A server Adapter can have multiple handlers
// A handler can have multiple tasks sets
// A task have a single external service adapter
package core

// ServerAdapter provides common interface for all the server resources.
// Example resources include:
// - REST API server
// - gRPC server
// - Kafka consumer
//
// It is intended to represent the individual interfaces on each exteranl service,
// not the services themselves; hence the name `ServiceAdapter`
type ServerAdapter interface {
    Serve() error
}

type ServerAdapterError struct {
    ServerAdapter ServerAdapter 
    Error error
}

type Handler struct {
	ID       string
	Name     string
	Protocol string
	Action   string
	TaskSets []TaskSet
}

type TaskSet struct {
	ServiceAdapter ServiceAdapter
	Concurrent             bool
}


