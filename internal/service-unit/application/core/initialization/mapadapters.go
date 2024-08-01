package initialization

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/primary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/resiliency"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	l "github.com/hanapedia/hexagon/pkg/operator/logger"
)

// mapSecondaryToPrimary map secondary adapter to primary adapter
func (su *ServiceUnit) mapSecondaryToPrimary() {
	for _, primaryConfig := range su.Config.AdapterConfigs {
		taskSet := su.newTaskSet(&primaryConfig)
		handler, err := su.newPrimaryAdapterHandler(&primaryConfig, taskSet)
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
func (su *ServiceUnit) newPrimaryAdapterHandler(primaryConfig *model.PrimaryAdapterSpec, taskSet []domain.TaskHandler) (domain.PrimaryAdapterHandler, error) {
	if primaryConfig.ServerConfig != nil {
		return domain.PrimaryAdapterHandler{
			ServiceName:  su.Name,
			ServerConfig: primaryConfig.ServerConfig,
			TaskSet:      taskSet,
		}, nil
	}
	if primaryConfig.ConsumerConfig != nil {
		return domain.PrimaryAdapterHandler{
			ServiceName:    su.Name,
			ConsumerConfig: primaryConfig.ConsumerConfig,
			TaskSet:        taskSet,
		}, nil
	}
	return domain.PrimaryAdapterHandler{}, errors.New("Failed to create primary adapter handler. No adapter config found.")
}

// newTaskSet creates task set from config
func (su *ServiceUnit) newTaskSet(primaryConfig *model.PrimaryAdapterSpec) []domain.TaskHandler {
	taskSet := make([]domain.TaskHandler, len(primaryConfig.TaskSpecs))
	for i, taskSpec := range primaryConfig.TaskSpecs {
		key := taskSpec.AdapterConfig.GetGroupByKey()
		client, ok := su.SecondaryAdapterClients[key]
		if !ok {
			l.Logger.Error("Client does not exist. ", "key=", key)
		}
		secondaryAdapter, err := secondary.NewSecondaryAdapter(taskSpec.AdapterConfig, client)
		if err != nil {
			l.Logger.Infof("Skipped interface: %s", err)
			continue
		}

		telemetryCtx := domain.NewTelemetryContext(su.Name, primaryConfig, taskSpec.AdapterConfig)

		taskSet[i] = resiliency.NewTaskHandler(telemetryCtx, taskSpec, secondaryAdapter)
	}

	return taskSet
}
