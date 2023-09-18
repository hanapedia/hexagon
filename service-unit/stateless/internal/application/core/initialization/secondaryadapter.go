package initialization

import (
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary"
)

// initializePrimaryAdapters prepare primary adapters
func (su *ServiceUnit) initializeSecondaryAdaptersClients() {
	for _, primaryConfig := range su.Config.AdapterConfigs {
		for _, step := range primaryConfig.Steps {
			key := step.AdapterConfig.GetGroupByKey()
			_, ok := (*su.SecondaryAdapters)[key]
			if !ok {
				(*su.SecondaryAdapters)[key] = secondary.NewSecondaryAdapterClient(step.AdapterConfig)
			}
		}
	}
}
