package initialization

import (
	"errors"
	"fmt"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/primary"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

type ServiceUnit struct {
	Name   string
	Config *model.ServiceUnitConfig
	// ServerAdapters hold the adapters for server processes from REST, gRPC
	ServerAdapters *map[constants.StatelessAdapterVariant]*ports.PrimaryPort
	// ConsumerAdapters hold the adapters for consumer processes from Kafka, RabbitMQ, Pular, etc
	ConsumerAdapters *map[string]*ports.PrimaryPort
	// SecondaryAdapters hold the persistent clients for secondary adapters
	SecondaryAdapters *map[string]ports.SecondaryAdapter
}

// NewServiceUnit initializes service unit object
func NewServiceUnit(serviceUnitConfig model.ServiceUnitConfig) ServiceUnit {
	serverAdapters := make(map[constants.StatelessAdapterVariant]*ports.PrimaryPort)
	consumerAdapters := make(map[string]*ports.PrimaryPort)

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
		logger.Logger.Infof("Serving '%s' server.", protocol)
		go func() {
			if err := (*serverAdapterCopy).Serve(); err != nil {
				errChan <- ports.PrimaryPortError{PrimaryPort: serverAdapterCopy, Error: err}
			}
		}()
	}

	for protocolAndAction, consumerAdapter := range *su.ConsumerAdapters {
		consumerAdapterCopy := consumerAdapter
		logger.Logger.Infof("Consumer '%s' started.", protocolAndAction)
		go func() {
			if err := (*consumerAdapterCopy).Serve(); err != nil {
				errChan <- ports.PrimaryPortError{PrimaryPort: consumerAdapterCopy, Error: err}
			}
		}()
	}
}

// Setup prepares primary adapters and maps secondary adapters to them
func (su *ServiceUnit) Setup() {
	su.initializePrimaryAdapters()
	su.mapSecondaryToPrimary()
}

// initializePrimaryAdapters prepare primary adapters
func (su *ServiceUnit) initializePrimaryAdapters() {
	for _, primaryConfig := range su.Config.IngressAdapterConfigs {
		if primaryConfig.StatelessIngressAdapterConfig != nil {
			su.initializeServerAdapter(*primaryConfig.StatelessIngressAdapterConfig)
			continue
		}
		if primaryConfig.BrokerIngressAdapterConfig != nil {
			su.initializeConsumerAdapter(*primaryConfig.BrokerIngressAdapterConfig)
			continue
		}
		logger.Logger.Fatal("Invalid primary adapter config.")
	}
}

// initializeServerAdapter prepare server adapters
func (su *ServiceUnit) initializeServerAdapter(config model.StatelessIngressAdapterConfig) {
	serverKey := getServerKey(config)
	_, ok := (*su.ServerAdapters)[serverKey]
	if !ok {
		(*su.ServerAdapters)[serverKey] = primary.NewServerAdapter(serverKey)
	}
}

// getServerKey retrieves server key from Stateless primary Adatper
func getServerKey(config model.StatelessIngressAdapterConfig) constants.StatelessAdapterVariant {
	return config.Variant
}

// initializeConsumerAdapter prepare consumer adapters
func (su *ServiceUnit) initializeConsumerAdapter(config model.BrokerIngressAdapterConfig) {
	consumerKey := getConsumerKey(config)
	_, ok := (*su.ConsumerAdapters)[consumerKey]
	if !ok {
		(*su.ConsumerAdapters)[consumerKey] = primary.NewConsumerAdapter(config.Variant, config.Topic)
	}
}

// getConsumerKey gets cosumer key from broker primary adapter
func getConsumerKey(config model.BrokerIngressAdapterConfig) string {
	return fmt.Sprintf("%s.%s", config.Variant, config.Topic)
}

// mapSecondaryToPrimary map secondary adapter to primary adapter
func (su *ServiceUnit) mapSecondaryToPrimary() {
	for _, primaryConfig := range su.Config.IngressAdapterConfigs {
		taskSet := su.createTaskSet(primaryConfig.Steps)
		handler, err := su.createPrimaryAdapterHandler(primaryConfig, taskSet)
		if err != nil {
			logger.Logger.Fatalf("Error creating handler: %v", err)
		}

		var primaryAdapter *ports.PrimaryPort
		if primaryConfig.StatelessIngressAdapterConfig != nil {
			primaryAdapter = (*su.ServerAdapters)[primaryConfig.StatelessIngressAdapterConfig.Variant]
		}
		if primaryConfig.BrokerIngressAdapterConfig != nil {
			consumerKey := getConsumerKey(*primaryConfig.BrokerIngressAdapterConfig)
			primaryAdapter = (*su.ConsumerAdapters)[consumerKey]
		}
		logger.Logger.Tracef("registering handler %s", handler.GetId(su.Name))

		err = primary.RegiserHandlerToPrimaryAdapter(su.Name, primaryAdapter, &handler)
		if err != nil {
			logger.Logger.Fatalf("Error registering handler to server adapter: %v", err)
		}
		logger.Logger.Infof("Successfully mapped '%s' handler", handler.GetId(su.Name))
	}
}

// createPrimaryAdapterHandler builds ingress adapter with given task set
func (su ServiceUnit) createPrimaryAdapterHandler(primaryConfig model.IngressAdapterSpec, taskSet *[]ports.TaskSet) (ports.PrimaryHandler, error) {
	if primaryConfig.StatelessIngressAdapterConfig != nil {
		return ports.PrimaryHandler{
			StatelessPrimaryAdapterConfig: primaryConfig.StatelessIngressAdapterConfig,
			TaskSets:                      *taskSet,
		}, nil
	}
	if primaryConfig.BrokerIngressAdapterConfig != nil {
		return ports.PrimaryHandler{
			BrokerPrimaryAdapterConfig: primaryConfig.BrokerIngressAdapterConfig,
			TaskSets:                   *taskSet,
		}, nil
	}
	return ports.PrimaryHandler{}, errors.New("Failed to create ingress adapter handler. No adapter config found.")
}

// createTaskSet creates task set from config
func (su ServiceUnit) createTaskSet(steps []model.Step) *[]ports.TaskSet {
	tasksets := make([]ports.TaskSet, len(steps))
	for i, step := range steps {
		secondaryAdapter, err := secondary.NewSecondaryAdapter(*step.EgressAdapterConfig, su.SecondaryAdapters)
		if err != nil {
			logger.Logger.Infof("Skipped interface: %s", err)
			continue
		}
		tasksets[i] = ports.TaskSet{SecondaryPort: secondaryAdapter, Concurrent: step.Concurrent}
	}

	return &tasksets
}
