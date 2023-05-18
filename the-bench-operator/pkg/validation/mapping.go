package validation

import (
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

// Validate adapter mapping across service unit configs
// serviceUnitConfigs: an array of serviceUnitConfig
func validateMapping(serviceUnitConfigs *[]model.ServiceUnitConfig) ConfigValidationError {
	serviceAdapterIds := mapIngressAdapters(serviceUnitConfigs)
	mappingErrors := mapEgressAdapters(serviceAdapterIds, serviceUnitConfigs)
	return ConfigValidationError{MappingErrors: mappingErrors}
}

// map ingress adapters to services
// returns the list of ids of ingress adapter
func mapIngressAdapters(serviceUnitConfigs *[]model.ServiceUnitConfig) []string {
	var serviceAdapterIds []string
	for _, serviceUnitConfig := range *serviceUnitConfigs {
		for _, ingressAdapterConfig := range serviceUnitConfig.IngressAdapterConfigs {
			serviceAdapterIds = append(serviceAdapterIds, ingressAdapterConfig.GetId(serviceUnitConfig.Name))
		}
	}
	return serviceAdapterIds
}

// map egress adapters to ingress adapters of services
// check if the id of egress adapter is found in the list of ids of ingress adapters
func mapEgressAdapters(serviceAdapterIds []string, serviceUnitConfigs *[]model.ServiceUnitConfig) []InvalidAdapterMappingError {
	var mappingErrors []InvalidAdapterMappingError
	for _, serviceUnitConfig := range *serviceUnitConfigs {
		for _, ingressAdapterConfig := range serviceUnitConfig.IngressAdapterConfigs {
			errs := mapAdapters(serviceAdapterIds, ingressAdapterConfig)
			if len(errs) != 0 {
				mappingErrors = append(mappingErrors, errs...)
				continue
			}
		}
	}
	return mappingErrors
}

// map egress adapters to ingress adapters of services
// conditionally handle the adapters
func mapAdapters(serviceAdapterIds []string, ingressAdapterConfig model.IngressAdapterSpec) []InvalidAdapterMappingError {
	var mappingErrors []InvalidAdapterMappingError
	for _, step := range ingressAdapterConfig.Steps {
		// ensure that egressAdapter is defined
		if step.EgressAdapterConfig == nil {
			continue
		}
		if step.EgressAdapterConfig.InternalEgressAdapterConfig != nil {
			continue
		}
		if ok := searchAdapterIds(serviceAdapterIds, *step.EgressAdapterConfig); !ok {
			mappingErrors = append(mappingErrors, NewInvalidEgressAdapterError(step.EgressAdapterConfig.GetId()))
		}
	}
	return mappingErrors
}

// perform the id search
func searchAdapterIds(serviceAdapterIds []string, egressAdapterConfig model.EgressAdapterConfig) bool {
	for _, id := range serviceAdapterIds {
		if id == egressAdapterConfig.GetId() {
			return true
		}
	}
	return false
}
