package model

import "log"

func ValidateServiceUnitConfigs(serviceUnitConfigs []ServiceUnitConfig) ConfigValidationError {
	// add missing fields
	addServieNameToAdapters(&serviceUnitConfigs)

	// validate service unit and adapters fields
	configValidationError := validateAdapterFields(serviceUnitConfigs)
	mappingErrors := validateAdapterMapping(serviceUnitConfigs)
	return ConfigValidationError{
		ServiceUnitFieldErrors: configValidationError.ServiceUnitFieldErrors,
		AdapterFieldErrors:     configValidationError.AdapterFieldErrors,
		MappingErrors:          mappingErrors,
	}
}

// validate all the field values
func validateAdapterFields(serviceUnitConfigs []ServiceUnitConfig) ConfigValidationError {
	var serviceUnitFieldErrors []InvalidServiceUnitFieldValueError
	var adapterFieldErrors []InvalidAdapterFieldValueError
	for _, serviceUnitConfig := range serviceUnitConfigs {
		configValidationError := ValidateServiceUnitConfigFields(serviceUnitConfig)
		serviceUnitFieldErrors = append(serviceUnitFieldErrors, configValidationError.ServiceUnitFieldErrors...)
		adapterFieldErrors = append(adapterFieldErrors, configValidationError.AdapterFieldErrors...)
	}
	return ConfigValidationError{ServiceUnitFieldErrors: serviceUnitFieldErrors, AdapterFieldErrors: adapterFieldErrors}
}

// validate fields for service config
func ValidateServiceUnitConfigFields(serviceUnitConfig ServiceUnitConfig) ConfigValidationError {
	var serviceUnitFieldErrors []InvalidServiceUnitFieldValueError
	serviceUnitFieldErrors = append(serviceUnitFieldErrors, serviceUnitConfig.Validate()...)

	var adapterFieldErrors []InvalidAdapterFieldValueError
	for i := range serviceUnitConfig.IngressAdapterConfigs {
		if serviceUnitConfig.IngressAdapterConfigs[i].StatelessIngressAdapterConfig != nil {
			serviceUnitConfig.IngressAdapterConfigs[i].StatelessIngressAdapterConfig.Service = serviceUnitConfig.Name
		}
		adapterFieldErrors = append(adapterFieldErrors, validateIngressAdapterConfig(&serviceUnitConfig.IngressAdapterConfigs[i])...)
		for _, step := range serviceUnitConfig.IngressAdapterConfigs[i].Steps {
			adapterFieldErrors = append(adapterFieldErrors, validateEgressAdapterConfig(step.EgressAdapterConfig)...)
		}
	}

	return ConfigValidationError{
		ServiceUnitFieldErrors: serviceUnitFieldErrors,
		AdapterFieldErrors:     adapterFieldErrors,
	}
}

func validateAdapterMapping(serviceUnitConfigs []ServiceUnitConfig) []InvalidAdapterMappingError {
	serviceAdapterIds := mapIngressAdapters(serviceUnitConfigs)
	mappingErrors := mapEgressAdapters(serviceAdapterIds, serviceUnitConfigs)
	return mappingErrors
}

// add service names to adapters if it does not exist
func addServieNameToAdapters(serviceUnitConfigs *[]ServiceUnitConfig) {
	for i := range *serviceUnitConfigs {
		for j := range (*serviceUnitConfigs)[i].IngressAdapterConfigs {
			if (*serviceUnitConfigs)[i].IngressAdapterConfigs[j].StatelessIngressAdapterConfig != nil {
				(*serviceUnitConfigs)[i].IngressAdapterConfigs[j].StatelessIngressAdapterConfig.Service = (*serviceUnitConfigs)[i].Name
			}
		}
	}
}

// map ingress adapters to services
func mapIngressAdapters(serviceUnitConfigs []ServiceUnitConfig) []string {
	var serviceAdapterIds []string
	for _, serviceUnitConfig := range serviceUnitConfigs {
		for _, ingressAdapterConfig := range serviceUnitConfig.IngressAdapterConfigs {
			serviceAdapterIds = append(serviceAdapterIds, generateIngressAdapterId(ingressAdapterConfig))
		}
	}
	return serviceAdapterIds
}

func generateIngressAdapterId(ingressAdapterConfig IngressAdapterConfig) string {
	var id string
	if ingressAdapterConfig.StatelessIngressAdapterConfig != nil {
		id = ingressAdapterConfig.StatelessIngressAdapterConfig.GetId()
	}
	if ingressAdapterConfig.BrokerIngressAdapterConfig != nil {
		id = ingressAdapterConfig.BrokerIngressAdapterConfig.GetId()
	}
	if ingressAdapterConfig.StatefulIngressAdapterConfig != nil {
		id = ingressAdapterConfig.StatefulIngressAdapterConfig.GetId()
	}
	return id
}

func generateEgressAdapterId(egressAdapterConfig EgressAdapterConfig) string {
	var id string
	if egressAdapterConfig.StatelessEgressAdapterConfig != nil {
		id = egressAdapterConfig.StatelessEgressAdapterConfig.GetId()
	}
	if egressAdapterConfig.BrokerEgressAdapterConfig != nil {
		id = egressAdapterConfig.BrokerEgressAdapterConfig.GetId()
	}
	if egressAdapterConfig.StatefulEgressAdapterConfig != nil {
		id = egressAdapterConfig.StatefulEgressAdapterConfig.GetId()
	}
	if egressAdapterConfig.InternalEgressAdapterConfig != nil {
		id = egressAdapterConfig.InternalEgressAdapterConfig.GetId()
	}
	return id
}

// Validate the fields of the ingress adapter configuration
func validateIngressAdapterConfig(ingressAdapterConfig *IngressAdapterConfig) []InvalidAdapterFieldValueError {
	var errs []InvalidAdapterFieldValueError
	if ingressAdapterConfig.StatelessIngressAdapterConfig != nil {
		errs = validateAdapter(ingressAdapterConfig.StatelessIngressAdapterConfig)
	}
	if ingressAdapterConfig.BrokerIngressAdapterConfig != nil {
		errs = validateAdapter(ingressAdapterConfig.BrokerIngressAdapterConfig)
	}
	if ingressAdapterConfig.StatefulIngressAdapterConfig != nil {
		if len(ingressAdapterConfig.Steps) > 0 {
			ingressAdapterConfig.Steps = []Step{} // makes sure that stateful service unit config have no steps defined
			log.Printf(
				"warning: unexpected steps definition on stateful ingress config for %s. These Steps will be ignored.",
				ingressAdapterConfig.StatefulIngressAdapterConfig.Name,
			)
		}
		errs = validateAdapter(ingressAdapterConfig.StatefulIngressAdapterConfig)
	}
	return errs
}

// Validate the fields of the egress adapter configuration
func validateEgressAdapterConfig(egressAdapterConfig EgressAdapterConfig) []InvalidAdapterFieldValueError {
	var errs []InvalidAdapterFieldValueError
	if egressAdapterConfig.StatelessEgressAdapterConfig != nil {
		errs = validateAdapter(*egressAdapterConfig.StatelessEgressAdapterConfig)
	}
	if egressAdapterConfig.BrokerEgressAdapterConfig != nil {
		errs = validateAdapter(*egressAdapterConfig.BrokerEgressAdapterConfig)
	}
	if egressAdapterConfig.StatefulEgressAdapterConfig != nil {
		errs = validateAdapter(*egressAdapterConfig.StatefulEgressAdapterConfig)
	}
	return errs
}

// map egress adapters to ingress adapters of services
func mapEgressAdapters(serviceAdapterIds []string, serviceUnitConfigs []ServiceUnitConfig) []InvalidAdapterMappingError {
	var mappingErrors []InvalidAdapterMappingError
	for _, serviceUnitConfig := range serviceUnitConfigs {
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
func mapAdapters(serviceAdapterIds []string, ingressAdapterConfig IngressAdapterConfig) []InvalidAdapterMappingError {
	var mappingErrors []InvalidAdapterMappingError
	for _, step := range ingressAdapterConfig.Steps {
		if step.EgressAdapterConfig.InternalEgressAdapterConfig != nil {
			continue
		}
		if ok := searchAdapterIds(serviceAdapterIds, step.EgressAdapterConfig); !ok {
			mappingErrors = append(mappingErrors, NewInvalidEgressAdapterError(generateEgressAdapterId(step.EgressAdapterConfig)))
		}
	}
	return mappingErrors
}

func searchAdapterIds(serviceAdapterIds []string, egressAdapterConfig EgressAdapterConfig) bool {
	for _, id := range serviceAdapterIds {
		if id == generateEgressAdapterId(egressAdapterConfig) {
			return true
		}
	}
	return false
}
