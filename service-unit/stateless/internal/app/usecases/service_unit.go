package usecases

import (
	"errors"
	"fmt"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	egressAdapterFactory "github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/factory"
	ingressAdapterFactory "github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/factory"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

type ServiceUnit struct {
	Name             string
	Config           *model.ServiceUnitConfig
	// ServerAdapters hold the adapters for server processes from REST, gRPC
	ServerAdapters   *map[constants.StatelessAdapterVariant]*core.IngressAdapter
	// ConsumerAdapters hold the adapters for consumer processes from Kafka, RabbitMQ, Pular, etc
	ConsumerAdapters *map[string]*core.IngressAdapter
	// EgressClients hold the persistent clients for egress adapters
	EgressClients    *map[string]core.EgressClient
}

// NewServiceUnit initializes service unit object
func NewServiceUnit(serviceUnitConfig model.ServiceUnitConfig) ServiceUnit {
	serverAdapters := make(map[constants.StatelessAdapterVariant]*core.IngressAdapter)
	consumerAdapters := make(map[string]*core.IngressAdapter)

	egressClients := make(map[string]core.EgressClient)

	return ServiceUnit{
		Name:             serviceUnitConfig.Name,
		Config:           &serviceUnitConfig,
		ServerAdapters:   &serverAdapters,
		ConsumerAdapters: &consumerAdapters,
		EgressClients:    &egressClients,
	}
}

// Start ingress adapters
func (su *ServiceUnit) Start(errChan chan core.IngressAdapterError) {
	for protocol, serverAdapter := range *su.ServerAdapters {
		serverAdapterCopy := serverAdapter
		logger.Logger.Infof("Serving '%s' server.", protocol)
		go func() {
			if err := (*serverAdapterCopy).Serve(); err != nil {
				errChan <- core.IngressAdapterError{IngressAdapter: serverAdapterCopy, Error: err}
			}
		}()
	}

	for protocolAndAction, consumerAdapter := range *su.ConsumerAdapters {
		consumerAdapterCopy := consumerAdapter
		logger.Logger.Infof("Consumer '%s' started.", protocolAndAction)
		go func() {
			if err := (*consumerAdapterCopy).Serve(); err != nil {
				errChan <- core.IngressAdapterError{IngressAdapter: consumerAdapterCopy, Error: err}
			}
		}()
	}
}

// Setup prepares ingress adapters and maps handlers to them
func (su *ServiceUnit) Setup() {
	su.initializeIngressAdapters()
	su.mapHandlersToIngressAdapters()
}

// initializeIngressAdapters prepare ingress adapters
func (su *ServiceUnit) initializeIngressAdapters() {
	for _, ingressAdapterConfig := range su.Config.IngressAdapterConfigs {
		if ingressAdapterConfig.StatelessIngressAdapterConfig != nil {
			su.initializeServerAdapter(*ingressAdapterConfig.StatelessIngressAdapterConfig)
			continue
		}
		if ingressAdapterConfig.BrokerIngressAdapterConfig != nil {
			su.initializeConsumerAdapter(*ingressAdapterConfig.BrokerIngressAdapterConfig)
			continue
		}
		logger.Logger.Fatal("Invalid ingress adapter config.")
	}
}

// initializeServerAdapter prepare server adapters
func (su *ServiceUnit) initializeServerAdapter(statelessAdapterConfig model.StatelessIngressAdapterConfig) {
	serverKey := getServerKey(statelessAdapterConfig)
	_, ok := (*su.ServerAdapters)[serverKey]
	if !ok {
		(*su.ServerAdapters)[serverKey] = ingressAdapterFactory.NewServerAdapter(serverKey)
	}
}

// getServerKey retrieves server key from Stateless Ingress Adatper
func getServerKey(statelessAdapterConfig model.StatelessIngressAdapterConfig) constants.StatelessAdapterVariant {
	return statelessAdapterConfig.Variant
}

// initializeConsumerAdapter prepare consumer adapters
func (su *ServiceUnit) initializeConsumerAdapter(brokerIngressAdapterConfig model.BrokerIngressAdapterConfig) {
	consumerKey := getConsumerKey(brokerIngressAdapterConfig)
	_, ok := (*su.ConsumerAdapters)[consumerKey]
	if !ok {
		(*su.ConsumerAdapters)[consumerKey] = ingressAdapterFactory.NewConsumerAdapter(brokerIngressAdapterConfig.Variant, brokerIngressAdapterConfig.Topic)
	}
}

// getConsumerKey gets cosumer key from broker ingress adapter
func getConsumerKey(brokerIngressAdapterConfig model.BrokerIngressAdapterConfig) string {
	return fmt.Sprintf("%s.%s", brokerIngressAdapterConfig.Variant, brokerIngressAdapterConfig.Topic)
}

// mapHandlersToIngressAdapters map egress adapter to ingress adapter
func (su *ServiceUnit) mapHandlersToIngressAdapters() {
	for _, ingressAdapterConfig := range su.Config.IngressAdapterConfigs {
		taskSets := su.mapTaskSet(ingressAdapterConfig.Steps)
		handler, err := su.createIngressAdapterHandler(ingressAdapterConfig, taskSets)
		if err != nil {
			logger.Logger.Fatalf("Error creating handler: %v", err)
		}

		var ingressAdapter *core.IngressAdapter
		if ingressAdapterConfig.StatelessIngressAdapterConfig != nil {
			ingressAdapter = (*su.ServerAdapters)[ingressAdapterConfig.StatelessIngressAdapterConfig.Variant]
		}
		if ingressAdapterConfig.BrokerIngressAdapterConfig != nil {
			consumerKey := getConsumerKey(*ingressAdapterConfig.BrokerIngressAdapterConfig)
			ingressAdapter = (*su.ConsumerAdapters)[consumerKey]
		}
		logger.Logger.Tracef("registering handler %s", handler.GetId(su.Name))

		err = ingressAdapterFactory.RegiserHandlerToIngressAdapter(su.Name, ingressAdapter, &handler)
		if err != nil {
			logger.Logger.Fatalf("Error registering handler to server adapter: %v", err)
		}
		logger.Logger.Infof("Successfully mapped '%s' handler", handler.GetId(su.Name))
	}
}

// createIngressAdapterHandler builds ingress adapter with given task set
func (su ServiceUnit) createIngressAdapterHandler(ingressAdapterConfig model.IngressAdapterSpec, taskSets *[]core.TaskSet) (core.IngressAdapterHandler, error) {
	if ingressAdapterConfig.StatelessIngressAdapterConfig != nil {
		return core.IngressAdapterHandler{
			StatelessIngressAdapterConfig: ingressAdapterConfig.StatelessIngressAdapterConfig,
			TaskSets:                      *taskSets,
		}, nil
	}
	if ingressAdapterConfig.BrokerIngressAdapterConfig != nil {
		return core.IngressAdapterHandler{
			BrokerIngressAdapterConfig: ingressAdapterConfig.BrokerIngressAdapterConfig,
			TaskSets:                   *taskSets,
		}, nil
	}
	return core.IngressAdapterHandler{}, errors.New("Failed to create ingress adapter handler. No adapter config found.")
}

// mapTaskSet creates task set from config
func (su ServiceUnit) mapTaskSet(steps []model.Step) *[]core.TaskSet {
	tasksets := make([]core.TaskSet, len(steps))
	for i, step := range steps {
		egressAdapter, err := egressAdapterFactory.NewEgressAdapter(*step.EgressAdapterConfig, su.EgressClients)
		if err != nil {
			logger.Logger.Infof("Skipped interface: %s", err)
			continue
		}
		tasksets[i] = core.TaskSet{EgressAdapter: egressAdapter, Concurrent: step.Concurrent}
	}

	return &tasksets
}
