package validation

import (
	"github.com/hanapedia/the-bench/config/logger"
	"github.com/hanapedia/the-bench/config/model"
)

func ValidateServiceUnitConfigs(serviceUnitConfigs *[]model.ServiceUnitConfig) model.ConfigValidationError {
	// validate service unit and adapters fields
	var configValidationError model.ConfigValidationError
	configValidationError.Extend(validateAdapterFields(serviceUnitConfigs))
	configValidationError.Extend(validateAdapterMapping(serviceUnitConfigs))
	return configValidationError
}

// validate all the field values
func validateAdapterFields(serviceUnitConfigs *[]model.ServiceUnitConfig) model.ConfigValidationError {
	var configValidationError model.ConfigValidationError
	for _, serviceUnitConfig := range *serviceUnitConfigs {
		configValidationError.Extend(ValidateServiceUnitConfigFields(&serviceUnitConfig))
	}
	return configValidationError
}

// validate fields for service config
func ValidateServiceUnitConfigFields(serviceUnitConfig *model.ServiceUnitConfig) model.ConfigValidationError {
	var configValidationError model.ConfigValidationError
	// use service name to stateful ingress adapters no matter it exist. Give warning when rewriting
	configValidationError.Extend(validateServiceUnitConfigFields(serviceUnitConfig))

	for i := range serviceUnitConfig.IngressAdapterConfigs {
		configValidationError.Extend(validateIngressAdapterConfig(serviceUnitConfig, &serviceUnitConfig.IngressAdapterConfigs[i]))
		for _, step := range serviceUnitConfig.IngressAdapterConfigs[i].Steps {
			// ensure that egressAdapter is defined
			if step.EgressAdapterConfig == nil {
				continue
			}
			configValidationError.Extend(validateEgressAdapterConfig(*step.EgressAdapterConfig))
		}
	}

	return configValidationError
}

func validateServiceUnitConfigFields(serviceUnitConfig *model.ServiceUnitConfig) model.ConfigValidationError {
	if len(serviceUnitConfig.IngressAdapterConfigs) > 1 {
		var statefulIngressAdapterConfig *model.StatefulIngressAdapterConfig
		for _, ingressAdapterConfig := range serviceUnitConfig.IngressAdapterConfigs {
			if ingressAdapterConfig.StatefulIngressAdapterConfig == nil {
				continue
			}
			statefulIngressAdapterConfig = ingressAdapterConfig.StatefulIngressAdapterConfig
			break
		}
		if statefulIngressAdapterConfig != nil {
			logger.Logger.Warnf("Stateful ingress adapter found, ignoring other ingress adapter definitions.\n")
			serviceUnitConfig.IngressAdapterConfigs = []model.IngressAdapterConfig{
				{
					StatefulIngressAdapterConfig: statefulIngressAdapterConfig,
				},
			}
		}
	}

	return model.ConfigValidationError{ServiceUnitFieldErrors: serviceUnitConfig.Validate()}
}

func validateAdapterMapping(serviceUnitConfigs *[]model.ServiceUnitConfig) model.ConfigValidationError {
	serviceAdapterIds := mapIngressAdapters(serviceUnitConfigs)
	mappingErrors := mapEgressAdapters(serviceAdapterIds, serviceUnitConfigs)
	return model.ConfigValidationError{MappingErrors: mappingErrors}
}

// map ingress adapters to services
func mapIngressAdapters(serviceUnitConfigs *[]model.ServiceUnitConfig) []string {
	var serviceAdapterIds []string
	for _, serviceUnitConfig := range *serviceUnitConfigs {
		for _, ingressAdapterConfig := range serviceUnitConfig.IngressAdapterConfigs {
			serviceAdapterIds = append(serviceAdapterIds, ingressAdapterConfig.GetId(serviceUnitConfig.Name))
		}
	}
	return serviceAdapterIds
}

// Validate the fields of the ingress adapter configuration
func validateIngressAdapterConfig(serviceUnitConfig *model.ServiceUnitConfig, ingressAdapterConfig *model.IngressAdapterConfig) model.ConfigValidationError {
	var adapterFieldErrors []model.InvalidAdapterFieldValueError
	if ingressAdapterConfig.StatelessIngressAdapterConfig != nil {
		adapterFieldErrors = model.ValidateIngressAdapter(serviceUnitConfig.Name, ingressAdapterConfig.StatelessIngressAdapterConfig)
	}
	if ingressAdapterConfig.BrokerIngressAdapterConfig != nil {
		adapterFieldErrors = model.ValidateIngressAdapter(serviceUnitConfig.Name, ingressAdapterConfig.BrokerIngressAdapterConfig)
	}
	if ingressAdapterConfig.StatefulIngressAdapterConfig != nil {
		if len(ingressAdapterConfig.Steps) > 0 {
			ingressAdapterConfig.Steps = []model.Step{} // makes sure that stateful service unit config have no steps defined
			logger.Logger.Warnf(
				"Steps definition found on stateful ingress config for %s. These Steps will be ignored.\n",
				serviceUnitConfig.Name,
			)
		}
		adapterFieldErrors = model.ValidateIngressAdapter(serviceUnitConfig.Name, ingressAdapterConfig.StatefulIngressAdapterConfig)
	}

	var stepFieldErrors []model.InvalidStepFieldValueError
	for _, step := range ingressAdapterConfig.Steps {
		stepFieldErrors = append(stepFieldErrors, step.Validate(serviceUnitConfig.Name, *ingressAdapterConfig)...)
	}
	return model.ConfigValidationError{AdapterFieldErrors: adapterFieldErrors, StepFieldErrors: stepFieldErrors}
}

// Validate the fields of the egress adapter configuration
func validateEgressAdapterConfig(egressAdapterConfig model.EgressAdapterConfig) model.ConfigValidationError {
	var errs []model.InvalidAdapterFieldValueError
	if egressAdapterConfig.StatelessEgressAdapterConfig != nil {
		errs = model.ValidateEgressAdapter(*egressAdapterConfig.StatelessEgressAdapterConfig)
	}
	if egressAdapterConfig.BrokerEgressAdapterConfig != nil {
		errs = model.ValidateEgressAdapter(*egressAdapterConfig.BrokerEgressAdapterConfig)
	}
	if egressAdapterConfig.StatefulEgressAdapterConfig != nil {
		errs = model.ValidateEgressAdapter(*egressAdapterConfig.StatefulEgressAdapterConfig)
	}
	return model.ConfigValidationError{AdapterFieldErrors: errs}
}

// map egress adapters to ingress adapters of services
func mapEgressAdapters(serviceAdapterIds []string, serviceUnitConfigs *[]model.ServiceUnitConfig) []model.InvalidAdapterMappingError {
	var mappingErrors []model.InvalidAdapterMappingError
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
func mapAdapters(serviceAdapterIds []string, ingressAdapterConfig model.IngressAdapterConfig) []model.InvalidAdapterMappingError {
	var mappingErrors []model.InvalidAdapterMappingError
	for _, step := range ingressAdapterConfig.Steps {
		// ensure that egressAdapter is defined
		if step.EgressAdapterConfig == nil {
			continue
		}
		if step.EgressAdapterConfig.InternalEgressAdapterConfig != nil {
			continue
		}
		if ok := searchAdapterIds(serviceAdapterIds, *step.EgressAdapterConfig); !ok {
			mappingErrors = append(mappingErrors, model.NewInvalidEgressAdapterError(step.EgressAdapterConfig.GetId()))
		}
	}
	return mappingErrors
}

func searchAdapterIds(serviceAdapterIds []string, egressAdapterConfig model.EgressAdapterConfig) bool {
	for _, id := range serviceAdapterIds {
		if id == egressAdapterConfig.GetId() {
			return true
		}
	}
	return false
}
