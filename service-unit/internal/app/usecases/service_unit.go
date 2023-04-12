package usecases

import (
	"fmt"
	"log"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	serverAdapterFactory "github.com/hanapedia/the-bench/service-unit/internal/infrastructure/server_adapter/factory"
	serviceAdapterFactory "github.com/hanapedia/the-bench/service-unit/internal/infrastructure/service_adapter/factory"
)

func NewServiceUnit(configLoader core.ConfigLoader) *core.ServiceUnit {
	config, err := configLoader.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	serverAdapters := mapServiceHandlerToServer(config.Name, &config.HandlerConfigs)

	serviceUnit := core.ServiceUnit{
		Name:           config.Name,
		ServerAdapters: serverAdapters,
	}

	return &serviceUnit
}

func mapServiceHandlerToServer(serviceName string, HandlerConfigs *[]core.HandlerConfig) *[]*core.ServerAdapter {
	serverAdapters := make([]*core.ServerAdapter, len(*HandlerConfigs))
	for _, HandlerConfig := range *HandlerConfigs {
		taskSets := mapTaskSet(&HandlerConfig.Flow)
		handler := core.Handler{
			Name:     HandlerConfig.Name,
			Protocol: HandlerConfig.Protocol,
			Action:   HandlerConfig.Action,
			ID: fmt.Sprintf(
				"%s.%s.%s.%s",
				serviceName,
				HandlerConfig.Protocol,
				HandlerConfig.Action,
				HandlerConfig.Name,
			),
			TaskSets: *taskSets,
		}

		err := serverAdapterFactory.UpsertServerAdapter(&serverAdapters, &handler)
		if err != nil {
			log.Fatalf("Error registering handler to server adapter: %v", err)
		}
	}
	return &serverAdapters
}

func mapTaskSet(steps *[]core.Step) *[]core.TaskSet {
	tasksets := make([]core.TaskSet, len(*steps))
	for i, step := range *steps {
		serviceAdapter, err := serviceAdapterFactory.NewServiceAdapterFromID(step.AdapterId)
		if err != nil {
			log.Printf("Skipped interface: %s", err)
			continue
		}
		tasksets[i] = core.TaskSet{ServiceAdapter: serviceAdapter, Concurrent: step.Concurrent}
	}

	return &tasksets
}
