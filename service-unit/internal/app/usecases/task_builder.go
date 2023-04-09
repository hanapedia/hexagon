package usecases

import (
	"fmt"
	"log"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	externalBuilder "github.com/hanapedia/the-bench/service-unit/internal/infrastructure/external_service_interface/builder"
)

type ServiceUnit struct {
	Name             string
	ServerInterfaces []ServerInterface
}

type ServerInterface struct {
	ID       string
	TaskSets []TaskSet
}

type TaskSet struct {
	ExternalServiceInterface core.ExternalServiceInterface
	Concurrent               bool
}

func (si *ServerInterface) Handle() {
	for _, ts := range si.TaskSets {
		ts.ExternalServiceInterface.Run()
	}
}

func NewServiceUnit(configLoader core.ConfigLoader) ServiceUnit {
	config, err := configLoader.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	serviceInterfaces := mapServerInterface(config.Name, &config.ServerInterfaceConfigs)

	serviceUnit := ServiceUnit{
		Name:             config.Name,
		ServerInterfaces: *serviceInterfaces,
	}
	return serviceUnit
}

func mapServerInterface(serviceName string, serverInterfaceConfigs *[]core.ServerInterfaceConfig) *[]ServerInterface {
	serverInterfaces := make([]ServerInterface, len(*serverInterfaceConfigs))
	for i, serverInterfaceConfig := range *serverInterfaceConfigs {
		taskSets := mapTaskSet(&serverInterfaceConfig.Flow)
		serverInterfaces[i] = ServerInterface{
			ID: fmt.Sprintf(
				"%s.%s.%s.%s",
				serviceName,
				serverInterfaceConfig.Protocol,
				serverInterfaceConfig.Action,
				serverInterfaceConfig.Name,
			),
			TaskSets: *taskSets,
		}
	}
	return &serverInterfaces
}

func mapTaskSet(steps *[]core.Step) *[]TaskSet {
	tasksets := make([]TaskSet, len(*steps))
	for i, step := range *steps {
		externalServiceInterface, err := externalBuilder.NewExternalServiceInterfaceFromID(step.InterfaceID)
		if err != nil {
			log.Printf("Skipped interface: %s", err)
			continue
		}
        tasksets[i] = TaskSet{ExternalServiceInterface: externalServiceInterface, Concurrent: step.Concurrent}
	}

	return &tasksets
}
