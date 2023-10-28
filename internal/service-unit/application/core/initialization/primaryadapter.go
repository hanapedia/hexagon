package initialization

import (
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	l "github.com/hanapedia/hexagon/pkg/operator/logger"
)

// initializePrimaryAdapters prepare primary adapters
func (su *ServiceUnit) initializePrimaryAdapters() {
	for _, primaryConfig := range su.Config.AdapterConfigs {
		if primaryConfig.ServerConfig != nil {
			su.initializeServerAdapter(primaryConfig.ServerConfig)
			continue
		}
		if primaryConfig.ConsumerConfig != nil {
			su.initializeConsumerAdapter(primaryConfig.ConsumerConfig)
			continue
		}
		l.Logger.Fatal("Invalid primary adapter config")
	}
}

// initializeServerAdapter prepare server adapters
func (su *ServiceUnit) initializeServerAdapter(config *model.ServerConfig) {
	serverKey := config.GetGroupByKey()
	_, ok := su.ServerAdapters[serverKey]
	if !ok {
		su.ServerAdapters[serverKey] = primary.NewServerAdapter(config)
	}
}

// initializeConsumerAdapter prepare consumer adapters
func (su *ServiceUnit) initializeConsumerAdapter(config *model.ConsumerConfig) {
	consumerKey := config.GetGroupByKey()
	_, ok := su.ConsumerAdapters[consumerKey]
	if !ok {
		su.ConsumerAdapters[consumerKey] = primary.NewConsumerAdapter(config)
	}
}

