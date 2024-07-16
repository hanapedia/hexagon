package initialization

import (
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary"
)

// initializePrimaryAdapters prepare primary adapters
func (su *ServiceUnit) initializeSecondaryAdaptersClients() {
	for _, primaryConfig := range su.Config.AdapterConfigs {
		for _, step := range primaryConfig.Tasks {
			if step.AdapterConfig.StressorConfig != nil {
				continue
			}
			key := step.AdapterConfig.GetGroupByKey()
			_, ok := su.SecondaryAdapterClients[key]
			if !ok {
				su.SecondaryAdapterClients[key] = secondary.NewSecondaryAdapterClient(step.AdapterConfig)
			}
		}
	}
}
