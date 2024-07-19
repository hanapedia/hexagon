package primary

import (
	"context"
	"sync"

	"github.com/hanapedia/hexagon/internal/service-unit/domain"
)

// PrimaryPort provides common interface for all the primary adapters.
// Example resources include:
// - REST API server
// - gRPC server
// - Kafka consumer
//
// It is intended to represent the individual interfaces on each exteranl service,
type PrimaryPort interface {
	// Serve starts primary port adapter with cancellable context and WaitGroup for graceful shutdown
	Serve(context.Context, *sync.WaitGroup) error
	// Register registers primary port handler to primary adapter instance
	Register(*domain.PrimaryHandler) error
}

type PrimaryPortError struct {
	PrimaryPort PrimaryPort
	Error       error
}

