package initialization

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/primary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	l "github.com/hanapedia/hexagon/pkg/operator/logger"
)

// mapSecondaryToPrimary map secondary adapter to primary adapter
func (su *ServiceUnit) mapSecondaryToPrimary() {
	for _, primaryConfig := range su.Config.AdapterConfigs {
		taskSet := su.newTaskSet(primaryConfig.Tasks)
		handler, err := su.newPrimaryAdapterHandler(primaryConfig, taskSet)
		if err != nil {
			l.Logger.Fatalf("Error creating handler: %v", err)
		}

		var primaryAdapter primary.PrimaryPort
		if primaryConfig.ServerConfig != nil {
			serverKey := primaryConfig.ServerConfig.GetGroupByKey()
			primaryAdapter = su.ServerAdapters[serverKey]
		}
		if primaryConfig.ConsumerConfig != nil {
			consumerKey := primaryConfig.ConsumerConfig.GetGroupByKey()
			primaryAdapter = su.ConsumerAdapters[consumerKey]
		}

		err = primaryAdapter.Register(&handler)
		if err != nil {
			l.Logger.Fatalf("Error registering handler to server adapter: %v", err)
		}
		l.Logger.Infof("Mapped '%s' handler", handler.GetId())
	}
}

// newPrimaryAdapterHandler builds primary adapter with given task set
func (su *ServiceUnit) newPrimaryAdapterHandler(primaryConfig model.PrimaryAdapterSpec, taskSet []domain.Task) (domain.PrimaryAdapterHandler, error) {
	if primaryConfig.ServerConfig != nil {
		return domain.PrimaryAdapterHandler{
			ServiceName: su.Name,
			ServerConfig: primaryConfig.ServerConfig,
			TaskSet:     taskSet,
		}, nil
	}
	if primaryConfig.ConsumerConfig != nil {
		return domain.PrimaryAdapterHandler{
			ServiceName: su.Name,
			ConsumerConfig: primaryConfig.ConsumerConfig,
			TaskSet:       taskSet,
		}, nil
	}
	return domain.PrimaryAdapterHandler{}, errors.New("Failed to create primary adapter handler. No adapter config found.")
}

// newTaskSet creates task set from config
func (su *ServiceUnit) newTaskSet(tasks []model.TaskSpec) []domain.Task {
	taskSet := make([]domain.Task, len(tasks))
	for i, task := range tasks {
		key := task.AdapterConfig.GetGroupByKey()
		client, ok := su.SecondaryAdapterClients[key]
		if !ok {
			l.Logger.Error("Client does not exist. ", "key=", key)
		}
		secondaryAdapter, err := secondary.NewSecondaryAdapter(task.AdapterConfig, client)
		if err != nil {
			l.Logger.Infof("Skipped interface: %s", err)
			continue
		}
		taskSet[i] = domain.Task{
			SecondaryPort: secondaryAdapter,
			Concurrent: task.Concurrent,
			OnError: task.AdapterConfig.OnError,
			TaskTimeout: task.AdapterConfig.TaskTimeout,
			CallTimeout: task.AdapterConfig.CallTimeout,
		}
	}

	return taskSet
}
