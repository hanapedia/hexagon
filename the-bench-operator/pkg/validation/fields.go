package validation

import (
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

// validate fields for a single service config
func ValidateFields(serviceUnitConfig *model.ServiceUnitConfig) ConfigValidationError {
	var configValidationError ConfigValidationError
	// use service name to stateful ingress adapters no matter it exist. Give warning when rewriting
	configValidationError.Extend(validateServiceUnitConfigFields(serviceUnitConfig))

	for i := range serviceUnitConfig.IngressAdapterConfigs {
		configValidationError.Extend(validateIngressAdapterConfigFields(serviceUnitConfig, &serviceUnitConfig.IngressAdapterConfigs[i]))
		for _, step := range serviceUnitConfig.IngressAdapterConfigs[i].Steps {
			// ensure that egressAdapter is defined
			if step.EgressAdapterConfig == nil {
				continue
			}
			configValidationError.Extend(validateEgressAdapterConfigFields(*step.EgressAdapterConfig))
		}
	}

	return configValidationError
}

// validate fields on service unit config
func validateServiceUnitConfigFields(serviceUnitConfig *model.ServiceUnitConfig) ConfigValidationError {
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
			serviceUnitConfig.IngressAdapterConfigs = []model.IngressAdapterSpec{
				{
					StatefulIngressAdapterConfig: statefulIngressAdapterConfig,
				},
			}
		}
	}
	return ValidateServiceUnitConfigFields(serviceUnitConfig)
}

// Validate the fields of the ingress adapter configuration
func validateIngressAdapterConfigFields(serviceUnitConfig *model.ServiceUnitConfig, ingressAdapterConfig *model.IngressAdapterSpec) ConfigValidationError {
	var adapterFieldErrors []InvalidAdapterFieldValueError
	if ingressAdapterConfig.StatelessIngressAdapterConfig != nil {
		adapterFieldErrors = ValidateStatelessIngressAdapterFields(
			ingressAdapterConfig.StatelessIngressAdapterConfig,
			serviceUnitConfig.Name,
		)
	}
	if ingressAdapterConfig.BrokerIngressAdapterConfig != nil {
		adapterFieldErrors = ValidateBrokerIngressAdapterFields(
			ingressAdapterConfig.BrokerIngressAdapterConfig,
			serviceUnitConfig.Name,
		)
	}
	if ingressAdapterConfig.StatefulIngressAdapterConfig != nil {
		if len(ingressAdapterConfig.Steps) > 0 {
			ingressAdapterConfig.Steps = []model.Step{} // makes sure that stateful service unit config have no steps defined
			logger.Logger.Warnf(
				"Steps definition found on stateful ingress config for %s. These Steps will be ignored.\n",
				serviceUnitConfig.Name,
			)
		}
		adapterFieldErrors = ValidateStatefulIngressAdapterFields(
			ingressAdapterConfig.StatefulIngressAdapterConfig,
			serviceUnitConfig.Name,
		)
	}
	var stepFieldErrors []InvalidStepFieldValueError
	for _, step := range ingressAdapterConfig.Steps {
		errs := ValidateStepFields(step, ingressAdapterConfig, serviceUnitConfig.Name)
		stepFieldErrors = append(stepFieldErrors, errs...)
	}
	return ConfigValidationError{AdapterFieldErrors: adapterFieldErrors, StepFieldErrors: stepFieldErrors}
}

// Validate the fields of the egress adapter configuration
func validateEgressAdapterConfigFields(egressAdapterConfig model.EgressAdapterConfig) ConfigValidationError {
	var errs []InvalidAdapterFieldValueError
	if egressAdapterConfig.StatelessEgressAdapterConfig != nil {
		errs = ValidateStatelessEgressAdapterFields(*egressAdapterConfig.StatelessEgressAdapterConfig)
	}
	if egressAdapterConfig.BrokerEgressAdapterConfig != nil {
		errs = ValidateBrokerEgressAdapterFields(*egressAdapterConfig.BrokerEgressAdapterConfig)
	}
	if egressAdapterConfig.StatefulEgressAdapterConfig != nil {
		errs = ValidateStatefulEgressAdapterFields(*egressAdapterConfig.StatefulEgressAdapterConfig)
	}
	return ConfigValidationError{AdapterFieldErrors: errs}
}
