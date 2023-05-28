package core

import (
	"context"
	"reflect"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

// EgressAdapter provides common interface for all the external service resouce.
// Example resources include:
// - REST API routes
// - gRPC methods
// - Kafka topic
// - Database table
//
// It is intended to represent the individual endpoints on each exteranl service,
// not the services themselves; hence the name `EgressAdapter`
type EgressAdapter interface {
	Call(context.Context) (string, error)
}

// Used to reuse clients to other serivces
// Wrapper interface, so the struct to implement this should have pointer to actual client
// types to implement this interface should have some sort of pointer to clients
type EgressClient interface {
	Close()
}

type EgressAdapterError struct {
	EgressAdapter *EgressAdapter
	Error         error
}

// LogEgressAdapterErrors logs the failed tasks
func LogEgressAdapterErrors(egressAdapterErrors *[]EgressAdapterError) {
	for _, egressAdapterError := range *egressAdapterErrors {
		logger.Logger.Errorf("Invocating %s failed: %s",
			reflect.TypeOf(egressAdapterError.EgressAdapter).Elem().Name(),
			egressAdapterError.Error,
		)
	}

}
