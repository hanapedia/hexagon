package initialization

import (
	"context"
	"sync"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/primary"
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

type ServiceUnit struct {
	Name   string
	Config *model.ServiceUnitConfig
	// ServerAdapters hold the adapters for server processes from REST, gRPC
	ServerAdapters map[string]primary.PrimaryPort
	// ConsumerAdapters hold the adapters for consumer processes from Kafka, RabbitMQ, Pular, etc
	ConsumerAdapters map[string]primary.PrimaryPort
	// SecondaryAdapterClients hold the persistent clients for secondary adapters
	SecondaryAdapterClients map[string]secondary.SecondaryAdapterClient
}

// NewServiceUnit initializes service unit object
func NewServiceUnit(serviceUnitConfig *model.ServiceUnitConfig) ServiceUnit {
	serverAdapters := make(map[string]primary.PrimaryPort)
	consumerAdapters := make(map[string]primary.PrimaryPort)

	secondaryAdapters := make(map[string]secondary.SecondaryAdapterClient)

	return ServiceUnit{
		Name:                    serviceUnitConfig.Name,
		Config:                  serviceUnitConfig,
		ServerAdapters:          serverAdapters,
		ConsumerAdapters:        consumerAdapters,
		SecondaryAdapterClients: secondaryAdapters,
	}
}

// Start primary adapters. it propagates context to primary adapters
func (su *ServiceUnit) Start(shutdownNotification context.Context, shutdownWaitGroup *sync.WaitGroup, errChan chan primary.PrimaryPortError) {
	for _, serverAdapter := range su.ServerAdapters {
		serverAdapterCopy := serverAdapter
		shutdownWaitGroup.Add(1)
		go func() {
			if err := serverAdapterCopy.Serve(shutdownNotification, shutdownWaitGroup); err != nil {
				errChan <- primary.PrimaryPortError{PrimaryPort: serverAdapterCopy, Error: err}
			}
		}()
	}

	for protocolAndAction, consumerAdapter := range su.ConsumerAdapters {
		consumerAdapterCopy := consumerAdapter
		logger.Logger.Infof("Consumer '%s' started.", protocolAndAction)
		shutdownWaitGroup.Add(1)
		go func() {
			if err := consumerAdapterCopy.Serve(shutdownNotification, shutdownWaitGroup); err != nil {
				errChan <- primary.PrimaryPortError{PrimaryPort: consumerAdapterCopy, Error: err}
			}
		}()
	}
}

// Setup prepares primary adapters and maps secondary adapters to them
func (su *ServiceUnit) Setup() {
	su.initializePrimaryAdapters()
	su.initializeSecondaryAdaptersClients()
	su.mapSecondaryToPrimary()
}

// Close closes all secondary adapter client connections
func (su *ServiceUnit) Close() {
	for key, client := range su.SecondaryAdapterClients	{
		logger.Logger.Infof("Closing %s client.", key)
		client.Close()
	}
}
