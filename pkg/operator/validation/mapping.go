package validation

import (
	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

// Validate adapter mapping across service unit configs
// serviceUnitConfigs: an array of serviceUnitConfig
func validateMapping(serviceUnitConfigs []model.ServiceUnitConfig) ConfigValidationError {
	serviceAdapterIds := mapPrimaryAdapters(serviceUnitConfigs)
	mappingErrors := mapSecondaryAdapters(serviceAdapterIds, serviceUnitConfigs)
	return ConfigValidationError{MappingErrors: mappingErrors}
}

// map primary adapters to services
// returns the list of ids of primary adapter
func mapPrimaryAdapters(serviceUnitConfigs []model.ServiceUnitConfig) []string {
	var serviceAdapterIds []string
	for _, serviceUnitConfig := range serviceUnitConfigs {
		for _, primaryAdapterConfig := range serviceUnitConfig.AdapterConfigs {
			serviceAdapterIds = append(serviceAdapterIds, primaryAdapterConfig.GetId(serviceUnitConfig.Name))
		}
	}
	return serviceAdapterIds
}

// map secondary adapters to primary adapters of services
// check if the id of secondary adapter is found in the list of ids of primary adapters
func mapSecondaryAdapters(serviceAdapterIds []string, serviceUnitConfigs []model.ServiceUnitConfig) []InvalidAdapterMappingError {
	var mappingErrors []InvalidAdapterMappingError
	for _, serviceUnitConfig := range serviceUnitConfigs {
		for _, primaryAdapterConfig := range serviceUnitConfig.AdapterConfigs {
			errs := mapAdapters(serviceAdapterIds, primaryAdapterConfig)
			if len(errs) != 0 {
				mappingErrors = append(mappingErrors, errs...)
				continue
			}
		}
	}
	return mappingErrors
}

// map secondary adapters to primary adapters of services
// conditionally handle the adapters
func mapAdapters(serviceAdapterIds []string, primaryAdapterConfig model.PrimaryAdapterSpec) []InvalidAdapterMappingError {
	var mappingErrors []InvalidAdapterMappingError
	for _, task := range primaryAdapterConfig.Tasks {
		// ensure that secondaryAdapter is defined
		if task.AdapterConfig == nil {
			continue
		}
		if ok := searchAdapterIds(serviceAdapterIds, *task.AdapterConfig); !ok {
			mappingErrors = append(mappingErrors, NewInvalidSecondaryAdapterError(task.AdapterConfig.GetId()))
		}
	}
	return mappingErrors
}

// perform the id search
func searchAdapterIds(serviceAdapterIds []string, secondaryAdapterConfig model.SecondaryAdapterConfig) bool {
	for _, id := range serviceAdapterIds {
		if id == secondaryAdapterConfig.GetId() {
			return true
		}
	}
	return false
}
