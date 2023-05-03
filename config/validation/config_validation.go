package validation

import (
	"github.com/hanapedia/the-bench/config/model"
)

func ValidateServiceUnitConfigs(serviceUnitConfigs *[]model.ServiceUnitConfig) model.ConfigValidationError {
	// validate service unit and adapters fields
	configValidationError := validateAdapterFields(serviceUnitConfigs)
	mappingErrors := validateAdapterMapping(serviceUnitConfigs)
	return model.ConfigValidationError{
		ServiceUnitFieldErrors: configValidationError.ServiceUnitFieldErrors,
		AdapterFieldErrors:     configValidationError.AdapterFieldErrors,
		MappingErrors:          mappingErrors,
	}
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
	var serviceUnitFieldErrors []model.InvalidServiceUnitFieldValueError
	// use service name to stateful ingress adapters no matter it exist. Give warning when rewriting
	serviceUnitFieldErrors = append(serviceUnitFieldErrors, validateServiceUnitConfigFields(serviceUnitConfig)...)

	var adapterFieldErrors []model.InvalidAdapterFieldValueError
	for i := range serviceUnitConfig.IngressAdapterConfigs {
		serviceUnitConfig.IngressAdapterConfigs[i] = validateServieNameOnAdapters(serviceUnitConfig.IngressAdapterConfigs[i], serviceUnitConfig.Name)

		adapterFieldErrors = append(adapterFieldErrors, validateIngressAdapterConfig(&serviceUnitConfig.IngressAdapterConfigs[i])...)
		for _, step := range serviceUnitConfig.IngressAdapterConfigs[i].Steps {
			adapterFieldErrors = append(adapterFieldErrors, validateEgressAdapterConfig(step.EgressAdapterConfig)...)
		}
	}

	return model.ConfigValidationError{
		ServiceUnitFieldErrors: serviceUnitFieldErrors,
		AdapterFieldErrors:     adapterFieldErrors,
	}
}

func validateServiceUnitConfigFields(serviceUnitConfig *model.ServiceUnitConfig) []model.InvalidServiceUnitFieldValueError {
	if len(serviceUnitConfig.IngressAdapterConfigs) > 1 {
		var statefulIngressAdapterConfig *model.StatefulAdapterConfig
		for _, ingressAdapterConfig := range serviceUnitConfig.IngressAdapterConfigs {
			if ingressAdapterConfig.StatefulIngressAdapterConfig == nil {
				continue
			}
			statefulIngressAdapterConfig = ingressAdapterConfig.StatefulIngressAdapterConfig
			break
		}
		if statefulIngressAdapterConfig != nil {
			logger.Warnf("Stateful ingress adapter found, ignoring other ingress adapter definitions.\n")
			serviceUnitConfig.IngressAdapterConfigs = []model.IngressAdapterConfig{
				{
					StatefulIngressAdapterConfig: statefulIngressAdapterConfig,
				},
			}
		}
	}

	return serviceUnitConfig.Validate()
}

// add service name to stateless ingress adapters if it does not exist
func validateServieNameOnAdapters(ingressAdapterConfig model.IngressAdapterConfig, serviceName string) model.IngressAdapterConfig {
	// ensure the service name consistecy for stateless ingress adapters
	if ingressAdapterConfig.StatelessIngressAdapterConfig != nil {
		addServieNameToStatelessAdapters(
			ingressAdapterConfig.StatelessIngressAdapterConfig,
			serviceName,
		)
	}
	// ensure the service name consistecy for stateful ingress adapters
	if ingressAdapterConfig.StatefulIngressAdapterConfig != nil {
		validateServiceNamOnStatefulAdapter(
			ingressAdapterConfig.StatefulIngressAdapterConfig,
			serviceName,
		)
	}
	return ingressAdapterConfig
}

// add service name to stateless ingress adapters if it does not exist
func addServieNameToStatelessAdapters(statelessIngressAdapterConfig *model.StatelessAdapterConfig, serviceName string) {
	if statelessIngressAdapterConfig.Service == "" {
		statelessIngressAdapterConfig.Service = serviceName
		logger.Infof(
			"Service field is undefined on stateless ingress adapter %s. Using Service Config service name.\n",
			statelessIngressAdapterConfig.GetId(),
		)
	} else if statelessIngressAdapterConfig.Service != serviceName {
		statelessIngressAdapterConfig.Service = serviceName
		logger.Warnf(
			"Service Config service name and ingress adapter does not match for ingress adapter %s. Resorting to Service Config service name for consistecy.\n",
			statelessIngressAdapterConfig.GetId(),
		)
	}
}

// ensure that the service name in Service unit config is identical to the service name in stateful ingress adapter
func validateServiceNamOnStatefulAdapter(statefulIngressAdapterConfig *model.StatefulAdapterConfig, serviceName string) {
	if statefulIngressAdapterConfig.Name != serviceName {
		statefulIngressAdapterConfig.Name = serviceName
		logger.Warnf(
			"Service Config service name and ingress adapter does not match for ingress adapter %s. Resorting to Service Config service name for consistecy.\n",
			statefulIngressAdapterConfig.GetId(),
		)
	}
}

func validateAdapterMapping(serviceUnitConfigs *[]model.ServiceUnitConfig) []model.InvalidAdapterMappingError {
	serviceAdapterIds := mapIngressAdapters(serviceUnitConfigs)
	mappingErrors := mapEgressAdapters(serviceAdapterIds, serviceUnitConfigs)
	return mappingErrors
}

// map ingress adapters to services
func mapIngressAdapters(serviceUnitConfigs *[]model.ServiceUnitConfig) []string {
	var serviceAdapterIds []string
	for _, serviceUnitConfig := range *serviceUnitConfigs {
		for _, ingressAdapterConfig := range serviceUnitConfig.IngressAdapterConfigs {
			serviceAdapterIds = append(serviceAdapterIds, generateIngressAdapterId(ingressAdapterConfig))
		}
	}
	return serviceAdapterIds
}

func generateIngressAdapterId(ingressAdapterConfig model.IngressAdapterConfig) string {
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

func generateEgressAdapterId(egressAdapterConfig model.EgressAdapterConfig) string {
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
func validateIngressAdapterConfig(ingressAdapterConfig *model.IngressAdapterConfig) []model.InvalidAdapterFieldValueError {
	var errs []model.InvalidAdapterFieldValueError
	if ingressAdapterConfig.StatelessIngressAdapterConfig != nil {
		errs = model.ValidateAdapter(ingressAdapterConfig.StatelessIngressAdapterConfig)
	}
	if ingressAdapterConfig.BrokerIngressAdapterConfig != nil {
		errs = model.ValidateAdapter(ingressAdapterConfig.BrokerIngressAdapterConfig)
	}
	if ingressAdapterConfig.StatefulIngressAdapterConfig != nil {
		if len(ingressAdapterConfig.Steps) > 0 {
			ingressAdapterConfig.Steps = []model.Step{} // makes sure that stateful service unit config have no steps defined
			logger.Warnf(
				"Steps definition found on stateful ingress config for %s. These Steps will be ignored.\n",
				ingressAdapterConfig.StatefulIngressAdapterConfig.Name,
			)
		}
		errs = model.ValidateAdapter(ingressAdapterConfig.StatefulIngressAdapterConfig)
	}
	return errs
}

// Validate the fields of the egress adapter configuration
func validateEgressAdapterConfig(egressAdapterConfig model.EgressAdapterConfig) []model.InvalidAdapterFieldValueError {
	var errs []model.InvalidAdapterFieldValueError
	if egressAdapterConfig.StatelessEgressAdapterConfig != nil {
		errs = model.ValidateAdapter(*egressAdapterConfig.StatelessEgressAdapterConfig)
	}
	if egressAdapterConfig.BrokerEgressAdapterConfig != nil {
		errs = model.ValidateAdapter(*egressAdapterConfig.BrokerEgressAdapterConfig)
	}
	if egressAdapterConfig.StatefulEgressAdapterConfig != nil {
		errs = model.ValidateAdapter(*egressAdapterConfig.StatefulEgressAdapterConfig)
	}
	return errs
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
		if step.EgressAdapterConfig.InternalEgressAdapterConfig != nil {
			continue
		}
		if ok := searchAdapterIds(serviceAdapterIds, step.EgressAdapterConfig); !ok {
			mappingErrors = append(mappingErrors, model.NewInvalidEgressAdapterError(generateEgressAdapterId(step.EgressAdapterConfig)))
		}
	}
	return mappingErrors
}

func searchAdapterIds(serviceAdapterIds []string, egressAdapterConfig model.EgressAdapterConfig) bool {
	for _, id := range serviceAdapterIds {
		if id == generateEgressAdapterId(egressAdapterConfig) {
			return true
		}
	}
	return false
}
