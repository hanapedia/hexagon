package usecases

import (
	"fmt"
	"log"

	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	egressAdapterFactory "github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/factory"
	ingressAdapterFactory "github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/factory"
)

// ServerAdapters hold the adapters for server processes from REST, gRPC
// ConsumerAdapters hold the adapters for consumer processes from Kafka, RabbitMQ, Pular, etc
// EgressConnections hold the persistent connections for egress adapters
type ServiceUnit struct {
	Name              string
	Config            *model.ServiceUnitConfig
	ServerAdapters    *map[constants.StatelessAdapterVariant]*core.IngressAdapter
	ConsumerAdapters  *map[string]*core.IngressAdapter
	EgressConnections *map[string]core.EgressConnection
}

func NewServiceUnit(serviceUnitConfig model.ServiceUnitConfig) ServiceUnit {
	serverAdapters := make(map[constants.StatelessAdapterVariant]*core.IngressAdapter)
	consumerAdapters := make(map[string]*core.IngressAdapter)

	egressConnections := make(map[string]core.EgressConnection)

	return ServiceUnit{Config: &serviceUnitConfig, ServerAdapters: &serverAdapters, ConsumerAdapters: &consumerAdapters, EgressConnections: &egressConnections}
}

// Start ingress adapters
func (su *ServiceUnit) Start(errChan chan core.IngressAdapterError) {
	for protocol, serverAdapter := range *su.ServerAdapters {
		serverAdapterCopy := serverAdapter
		log.Printf("Serving '%s' server.", protocol)
		go func() {
			if err := (*serverAdapterCopy).Serve(); err != nil {
				errChan <- core.IngressAdapterError{IngressAdapter: serverAdapterCopy, Error: err}
			}
		}()
	}

	for protocolAndAction, consumerAdapter := range *su.ConsumerAdapters {
		consumerAdapterCopy := consumerAdapter
		log.Printf("Consumer '%s' started.", protocolAndAction)
		go func() {
			if err := (*consumerAdapterCopy).Serve(); err != nil {
				errChan <- core.IngressAdapterError{IngressAdapter: consumerAdapterCopy, Error: err}
			}
		}()
	}
}

// Prepares ingress adapters and maps handlers to them
func (su *ServiceUnit) Setup() {
	su.initializeIngressAdapters()
	su.mapHandlersToIngressAdapters()
}

// Prepare ingress adapters
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
		log.Fatal("Invalid ingress adapter config.")
	}
}

// Prepare server adapters
func (su *ServiceUnit) initializeServerAdapter(statelessAdapterConfig model.StatelessAdapterConfig) {
	serverKey := getServerKey(statelessAdapterConfig)
	_, ok := (*su.ServerAdapters)[serverKey]
	if !ok {
		(*su.ServerAdapters)[serverKey] = ingressAdapterFactory.NewServerAdapter(serverKey)
	}
}

// get server key
func getServerKey(statelessAdapterConfig model.StatelessAdapterConfig) constants.StatelessAdapterVariant {
	return statelessAdapterConfig.Variant
}

// Prepare consumer adapters
func (su *ServiceUnit) initializeConsumerAdapter(brokerIngressAdapterConfig model.BrokerAdapterConfig) {
	consumerKey := getConsumerKey(brokerIngressAdapterConfig)
	_, ok := (*su.ConsumerAdapters)[consumerKey]
	if !ok {
		(*su.ConsumerAdapters)[consumerKey] = ingressAdapterFactory.NewConsumerAdapter(brokerIngressAdapterConfig.Variant, brokerIngressAdapterConfig.Topic)
	}
}

func getConsumerKey(brokerIngressAdapterConfig model.BrokerAdapterConfig) string {
	return fmt.Sprintf("%s.%s", brokerIngressAdapterConfig.Variant, brokerIngressAdapterConfig.Topic)
}

// Handler configuration can omit the serivce name so the service name needs to be assigned for better logging
func assignServicename(service string, statelessAdapterConfig *model.StatelessAdapterConfig) *model.StatelessAdapterConfig {
	if statelessAdapterConfig.Service == "" {
		statelessAdapterConfig.Service = service
	}
	return statelessAdapterConfig
}

// Map handlers to ingress adapters
func (su *ServiceUnit) mapHandlersToIngressAdapters() {
	for _, ingressAdapterConfig := range su.Config.IngressAdapterConfigs {
		taskSets := su.mapTaskSet(ingressAdapterConfig.Steps)
		handler := core.IngressAdapterHandler{
			StatelessIngressAdapterConfig: assignServicename(su.Name, ingressAdapterConfig.StatelessIngressAdapterConfig),
			BrokerIngressAdapterConfig:    ingressAdapterConfig.BrokerIngressAdapterConfig,
			TaskSets:                      *taskSets,
		}

		var ingressAdapter *core.IngressAdapter
		if ingressAdapterConfig.StatelessIngressAdapterConfig != nil {
			ingressAdapter = (*su.ServerAdapters)[ingressAdapterConfig.StatelessIngressAdapterConfig.Variant]
		}
		if ingressAdapterConfig.BrokerIngressAdapterConfig != nil {
			consumerKey := getConsumerKey(*ingressAdapterConfig.BrokerIngressAdapterConfig)
			ingressAdapter = (*su.ConsumerAdapters)[consumerKey]
		}
		log.Printf("registering handler %s", handler.GetId())

		err := ingressAdapterFactory.RegiserHandlerToIngressAdapter(ingressAdapter, &handler)
		if err != nil {
			log.Fatalf("Error registering handler to server adapter: %v", err)
		}
		log.Printf("Successfully mapped '%s' handler.", handler.GetId())
	}
}

// Create task set from config
func (su ServiceUnit) mapTaskSet(steps []model.Step) *[]core.TaskSet {
	tasksets := make([]core.TaskSet, len(steps))
	for i, step := range steps {
		egressAdapter, err := egressAdapterFactory.NewEgressAdapter(step.EgressAdapterConfig, su.EgressConnections)
		if err != nil {
			log.Printf("Skipped interface: %s", err)
			continue
		}
		tasksets[i] = core.TaskSet{EgressAdapter: egressAdapter, Concurrent: step.Concurrent}
	}

	return &tasksets
}
