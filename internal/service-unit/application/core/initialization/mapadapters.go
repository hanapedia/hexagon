package initialization

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	l "github.com/hanapedia/hexagon/pkg/operator/logger"
)

// mapSecondaryToPrimary map secondary adapter to primary adapter
func (su *ServiceUnit) mapSecondaryToPrimary() {
	for _, primaryConfig := range su.Config.AdapterConfigs {
		taskSet := su.newTaskSet(primaryConfig.Steps)
		handler, err := su.newPrimaryAdapterHandler(primaryConfig, taskSet)
		if err != nil {
			l.Logger.Fatalf("Error creating handler: %v", err)
		}

		var primaryAdapter ports.PrimaryPort
		if primaryConfig.ServerConfig != nil {
			serverKey := primaryConfig.ServerConfig.GetGroupByKey()
			primaryAdapter = su.ServerAdapters[serverKey]
		}
		if primaryConfig.ConsumerConfig != nil {
			consumerKey := primaryConfig.ConsumerConfig.GetGroupByKey()
			primaryAdapter = su.ConsumerAdapters[consumerKey]
		}

		err = primary.RegiserHandlerToPrimaryAdapter(su.Name, primaryAdapter, &handler)
		if err != nil {
			l.Logger.Fatalf("Error registering handler to server adapter: %v", err)
		}
		l.Logger.Infof("Mapped '%s' handler", handler.GetId())
	}
}

// newPrimaryAdapterHandler builds primary adapter with given task set
func (su *ServiceUnit) newPrimaryAdapterHandler(primaryConfig model.PrimaryAdapterSpec, taskSet []ports.Task) (ports.PrimaryHandler, error) {
	if primaryConfig.ServerConfig != nil {
		return ports.PrimaryHandler{
			ServiceName: su.Name,
			ServerConfig: primaryConfig.ServerConfig,
			TaskSet:     taskSet,
		}, nil
	}
	if primaryConfig.ConsumerConfig != nil {
		return ports.PrimaryHandler{
			ServiceName: su.Name,
			ConsumerConfig: primaryConfig.ConsumerConfig,
			TaskSet:       taskSet,
		}, nil
	}
	return ports.PrimaryHandler{}, errors.New("Failed to create primary adapter handler. No adapter config found.")
}

// newTaskSet creates task set from config
func (su *ServiceUnit) newTaskSet(steps []model.Step) []ports.Task {
	taskSet := make([]ports.Task, len(steps))
	for i, step := range steps {
		key := step.AdapterConfig.GetGroupByKey()
		client, ok := su.SecondaryAdapterClients[key]
		if !ok {
			l.Logger.Error("Client does not exist. ", "key=", key)
		}
		secondaryAdapter, err := secondary.NewSecondaryAdapter(step.AdapterConfig, client)
		if err != nil {
			l.Logger.Infof("Skipped interface: %s", err)
			continue
		}
		taskSet[i] = ports.Task{SecondaryPort: secondaryAdapter, Concurrent: step.Concurrent}
	}

	return taskSet
}
