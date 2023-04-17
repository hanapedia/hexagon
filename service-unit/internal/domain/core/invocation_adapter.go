package core

// InvocationAdapter provides common interface for all the external service resouce.
// Example resources include:
// - REST API routes
// - gRPC methods
// - Kafka topic
// - Database table
//
// It is intended to represent the individual endpoints on each exteranl service,
// not the services themselves; hence the name `InvocationAdapter`
type InvocationAdapter interface {
	Call() (string, error)
}

// Used to reuse connections to other serivces
// Wrapper interface, so the struct to implement this should have pointer to actual connection
type Connection interface {
	Close()
}

type InvocationAdapterError struct {
	InvocationAdapter *InvocationAdapter
	Error             error
}
