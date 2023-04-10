package usecases

import (
	"fmt"
	"log"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	externalBuilder "github.com/hanapedia/the-bench/service-unit/internal/infrastructure/external_service_adapter/builder"
)

type ServiceUnit struct {
	Name            string
	serverAdapters  []core.ServerAdapter
	ServiceHandlers []ServiceHandler
}

type ServiceHandler struct {
	ID       string
	Name     string
	Protocol string
	Action   string
	TaskSets []TaskSet
}

type TaskSet struct {
	ExternalServiceAdapter core.ExternalServiceAdapter
	Concurrent             bool
}

func (si *ServiceHandler) Handle() {
	for _, ts := range si.TaskSets {
		ts.ExternalServiceAdapter.Call()
	}
}

func NewServiceUnit(configLoader core.ConfigLoader) ServiceUnit {
	config, err := configLoader.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	serviceHandler := mapServiceHandler(config.Name, &config.ServerInterfaceConfigs)

	serviceUnit := ServiceUnit{
		Name:            config.Name,
		ServiceHandlers: *serviceHandler,
	}
	return serviceUnit
}

func mapServiceHandler(serviceName string, serviceHandlerConfigs *[]core.ServiceHandlerConfig) *[]ServiceHandler {
	serviceHandler := make([]ServiceHandler, len(*serviceHandlerConfigs))
	for i, serviceHandlerConfig := range *serviceHandlerConfigs {
		taskSets := mapTaskSet(&serviceHandlerConfig.Flow)
		serviceHandler[i] = ServiceHandler{
			Name:     serviceHandlerConfig.Name,
			Protocol: serviceHandlerConfig.Protocol,
			Action:   serviceHandlerConfig.Action,
			ID: fmt.Sprintf(
				"%s.%s.%s.%s",
				serviceName,
				serviceHandlerConfig.Protocol,
				serviceHandlerConfig.Action,
				serviceHandlerConfig.Name,
			),
			TaskSets: *taskSets,
		}
	}
	return &serviceHandler
}

func mapTaskSet(steps *[]core.Step) *[]TaskSet {
	tasksets := make([]TaskSet, len(*steps))
	for i, step := range *steps {
		externalServiceAdapter, err := externalBuilder.NewExternalServiceAdapterFromID(step.AdapterId)
		if err != nil {
			log.Printf("Skipped interface: %s", err)
			continue
		}
		tasksets[i] = TaskSet{ExternalServiceAdapter: externalServiceAdapter, Concurrent: step.Concurrent}
	}

	return &tasksets
}
