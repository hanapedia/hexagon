package initialization

import (
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

// initializePrimaryAdapters prepare primary adapters
func (su *ServiceUnit) initializeSecondaryAdaptersClients() {
	for _, primaryConfig := range su.Config.AdapterConfigs {
		for _, task := range primaryConfig.TaskSpecs {
			key := task.AdapterConfig.GetGroupByKey()
			_, ok := su.SecondaryAdapterClients[key]
			if !ok {
				logger.Logger.WithField("key", key).Infof("creating new client")
				su.SecondaryAdapterClients[key] = secondary.NewSecondaryAdapterClient(task.AdapterConfig)
			}
		}
	}
}
