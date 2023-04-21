package usecases

import (
	"fmt"
	"log"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	egressAdapterFactory "github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/factory"
	ingressAdapterFactory "github.com/hanapedia/the-bench/service-unit/internal/infrastructure/ingress/factory"
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
)

// ServerAdapters hold the adapters for server processes from REST, gRPC
// ConsumerAdapters hold the adapters for consumer processes from Kafka, RabbitMQ, Pular, etc
// EgressConnections hold the persistent connections for egress adapters
type ServiceUnit struct {
	Name              string
	Config            *core.ServiceUnitConfig
	ServerAdapters    *map[constants.AdapterProtocol]*core.IngressAdapter
	ConsumerAdapters  *map[string]*core.IngressAdapter
	EgressConnections *map[string]core.EgressConnection
}

func NewServiceUnit(format string) ServiceUnit {
	configLoader := NewConfigLoader(format)
	config, err := configLoader.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	serverAdapters := make(map[constants.AdapterProtocol]*core.IngressAdapter)
	consumerAdapters := make(map[string]*core.IngressAdapter)

	egressConnections := make(map[string]core.EgressConnection)

	return ServiceUnit{Config: &config, ServerAdapters: &serverAdapters, ConsumerAdapters: &consumerAdapters, EgressConnections: &egressConnections}
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
	for _, handlerConfig := range su.Config.HandlerConfigs {
		switch handlerConfig.Protocol {
		case constants.REST:
			su.initializeServerAdapter(handlerConfig.Protocol)
		case constants.KAFKA:
			su.initializeConsumerAdapter(handlerConfig.Protocol, handlerConfig.Action)
		}
	}
}

// Prepare server adapters
func (su *ServiceUnit) initializeServerAdapter(protocol constants.AdapterProtocol) {
	_, ok := (*su.ServerAdapters)[protocol]
	if !ok {
		(*su.ServerAdapters)[protocol] = ingressAdapterFactory.NewServerAdapter(protocol)
	}
}

// Prepare consumer adapters
func (su *ServiceUnit) initializeConsumerAdapter(protocol constants.AdapterProtocol, action string) {
	consumerKey := fmt.Sprintf("%s.%s", protocol, action)
	_, ok := (*su.ConsumerAdapters)[consumerKey]
	if !ok {
		(*su.ConsumerAdapters)[consumerKey] = ingressAdapterFactory.NewConsumerAdapter(protocol, action)
	}
}

// Map handlers to ingress adapters
func (su *ServiceUnit) mapHandlersToIngressAdapters() {
	for _, handlerConfig := range su.Config.HandlerConfigs {
        log.Printf("start mapping %s", handlerConfig.Name)
		taskSets := su.mapTaskSet(handlerConfig.Steps)
		handler := core.Handler{
			Name:     handlerConfig.Name,
			Protocol: handlerConfig.Protocol,
			Action:   handlerConfig.Action,
			ID: fmt.Sprintf(
				"%s.%s.%s.%s",
				su.Name,
				handlerConfig.Protocol,
				handlerConfig.Action,
				handlerConfig.Name,
			),
			TaskSets: *taskSets,
		}
        log.Printf("taskset generated %s", handlerConfig.Name)

        var ingressAdapter *core.IngressAdapter
		switch handlerConfig.Protocol {
		case constants.REST:
            ingressAdapter = (*su.ServerAdapters)[handler.Protocol]
		case constants.KAFKA:
            consumerKey := fmt.Sprintf("%s.%s", handler.Protocol, handler.Action)
            ingressAdapter = (*su.ConsumerAdapters)[consumerKey]
		}
        log.Printf("registering handler %s", handlerConfig.Name)

		err := ingressAdapterFactory.RegiserHandlerToIngressAdapter(ingressAdapter, &handler)
		if err != nil {
			log.Fatalf("Error registering handler to server adapter: %v", err)
		}
		log.Printf("Successfully mapped '%s' handler to '%s' server", handler.Name, handler.Protocol)
	}
}

// Create task set from config
func (su ServiceUnit) mapTaskSet(steps []core.Step) *[]core.TaskSet {
	tasksets := make([]core.TaskSet, len(steps))
	for i, step := range steps {
		egressAdapter, err := egressAdapterFactory.NewEgressAdapter(step.AdapterId, su.EgressConnections)
		if err != nil {
			log.Printf("Skipped interface: %s", err)
			continue
		}
		tasksets[i] = core.TaskSet{EgressAdapter: egressAdapter, Concurrent: step.Concurrent}
	}

	return &tasksets
}

