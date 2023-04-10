package core

// ExternalServiceAdapter provides common interface for all the external service resouce.
// Example resources include:
// - REST API routes
// - gRPC methods
// - Kafka topic
// - Database table
//
// It is intended to represent the individual interfaces on each exteranl service,
// not the services themselves; hence the name `ExternalServiceAdapter`
type ExternalServiceAdapter interface {
	Call() (string, error)
}
