package usecases

import (
	"fmt"
	"log"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	serverAdapterFactory "github.com/hanapedia/the-bench/service-unit/internal/infrastructure/server_adapter/factory"
	invocationAdapterFactory "github.com/hanapedia/the-bench/service-unit/internal/infrastructure/invocation_adapter/factory"
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
)

func NewServiceUnit(configLoader core.ConfigLoader) core.ServiceUnit {
	config, err := configLoader.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
    log.Println("Config successfully loaded.")

	serverAdapters := mapServiceHandlerToServer(config.Name, &config.HandlerConfigs)
    log.Println("Handler to Server mapping completed.")

	serviceUnit := core.ServiceUnit{
		Name:           config.Name,
		ServerAdapters: serverAdapters,
	}

	return serviceUnit
}

func mapServiceHandlerToServer(serviceName string, HandlerConfigs *[]core.HandlerConfig) map[constants.ServerAdapterProtocol]*core.ServerAdapter {
	serverAdapters := make(map[constants.ServerAdapterProtocol]*core.ServerAdapter)
	for _, handlerConfig := range *HandlerConfigs {
		taskSets := mapTaskSet(&handlerConfig.Flow)
		handler := core.Handler{
			Name:     handlerConfig.Name,
			Protocol: handlerConfig.Protocol,
			Action:   handlerConfig.Action,
			ID: fmt.Sprintf(
				"%s.%s.%s.%s",
				serviceName,
				handlerConfig.Protocol,
				handlerConfig.Action,
				handlerConfig.Name,
			),
			TaskSets: *taskSets,
		}
		serverAdapterProtocol := constants.ServerAdapterProtocol(handler.Protocol)
		_, ok := serverAdapters[serverAdapterProtocol]
		if !ok {
			serverAdapters[serverAdapterProtocol] = serverAdapterFactory.NewServerAdapter(serverAdapterProtocol)
		}
		err := serverAdapterFactory.RegiserHandlerToServerAdapter(serverAdapterProtocol, serverAdapters[serverAdapterProtocol], &handler)
		if err != nil {
			log.Fatalf("Error registering handler to server adapter: %v", err)
		}
        log.Printf("Successfully mapped '%s' handler to '%s' server", handler.Name, serverAdapterProtocol)
	}
	return serverAdapters
}

func mapTaskSet(steps *[]core.Step) *[]core.TaskSet {
	tasksets := make([]core.TaskSet, len(*steps))
	for i, step := range *steps {
		invocationAdapter, err := invocationAdapterFactory.NewInvocationAdapterFromID(step.AdapterId)
		if err != nil {
			log.Printf("Skipped interface: %s", err)
			continue
		}
		tasksets[i] = core.TaskSet{InvocationAdapter: invocationAdapter, Concurrent: step.Concurrent}
	}

	return &tasksets
}
