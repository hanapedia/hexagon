package initialization

import (
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	l "github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

type ServiceUnit struct {
	Name   string
	Config *model.ServiceUnitConfig
	// ServerAdapters hold the adapters for server processes from REST, gRPC
	ServerAdapters *map[string]ports.PrimaryPort
	// ConsumerAdapters hold the adapters for consumer processes from Kafka, RabbitMQ, Pular, etc
	ConsumerAdapters *map[string]ports.PrimaryPort
	// SecondaryAdapters hold the persistent clients for secondary adapters
	SecondaryAdapters *map[string]ports.SecondaryAdapter
}

// NewServiceUnit initializes service unit object
func NewServiceUnit(serviceUnitConfig model.ServiceUnitConfig) ServiceUnit {
	serverAdapters := make(map[string]ports.PrimaryPort)
	consumerAdapters := make(map[string]ports.PrimaryPort)

	secondaryAdapters := make(map[string]ports.SecondaryAdapter)

	return ServiceUnit{
		Name:              serviceUnitConfig.Name,
		Config:            &serviceUnitConfig,
		ServerAdapters:    &serverAdapters,
		ConsumerAdapters:  &consumerAdapters,
		SecondaryAdapters: &secondaryAdapters,
	}
}

// Start primary adapters
func (su *ServiceUnit) Start(errChan chan ports.PrimaryPortError) {
	for protocol, serverAdapter := range *su.ServerAdapters {
		serverAdapterCopy := serverAdapter
		l.Logger.Infof("Serving '%s' server.", protocol)
		go func() {
			if err := serverAdapterCopy.Serve(); err != nil {
				errChan <- ports.PrimaryPortError{PrimaryPort: serverAdapterCopy, Error: err}
			}
		}()
	}

	for protocolAndAction, consumerAdapter := range *su.ConsumerAdapters {
		consumerAdapterCopy := consumerAdapter
		l.Logger.Infof("Consumer '%s' started.", protocolAndAction)
		go func() {
			if err := consumerAdapterCopy.Serve(); err != nil {
				errChan <- ports.PrimaryPortError{PrimaryPort: consumerAdapterCopy, Error: err}
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
