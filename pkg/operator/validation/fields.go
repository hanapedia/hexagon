package validation

import (
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

// validate fields for a single service config
func ValidateFields(serviceUnitConfig *model.ServiceUnitConfig) ConfigValidationError {
	var configValidationError ConfigValidationError
	// use service name to repository adapters no matter it exist. Give warning when rewriting
	configValidationError.Extend(validateServiceUnitConfigFields(serviceUnitConfig))

	for i := range serviceUnitConfig.AdapterConfigs {
		configValidationError.Extend(validatePrimaryAdapterConfigFields(serviceUnitConfig, &serviceUnitConfig.AdapterConfigs[i]))
		for _, step := range serviceUnitConfig.AdapterConfigs[i].Steps {
			// ensure that secondary adapter is defined
			if step.AdapterConfig == nil {
				continue
			}
			configValidationError.Extend(validateSecondaryAdapterConfigFields(*step.AdapterConfig))
		}
	}

	return configValidationError
}

// validate fields on service unit config
func validateServiceUnitConfigFields(serviceUnitConfig *model.ServiceUnitConfig) ConfigValidationError {
	if len(serviceUnitConfig.AdapterConfigs) > 1 {
		var repositoryConfig *model.RepositoryConfig
		for _, primaryAdapterConfig := range serviceUnitConfig.AdapterConfigs {
			if primaryAdapterConfig.RepositoryConfig == nil {
				continue
			}
			repositoryConfig = primaryAdapterConfig.RepositoryConfig
			break
		}
		if repositoryConfig != nil {
			logger.Logger.Warnf("repository adapter found, ignoring other primary adapter definitions.\n")
			serviceUnitConfig.AdapterConfigs = []model.PrimaryAdapterSpec{
				{
					RepositoryConfig: repositoryConfig,
				},
			}
		}
	}
	return ValidateServiceUnitConfigFields(serviceUnitConfig)
}

// Validate the fields of the primary adapter configuration
func validatePrimaryAdapterConfigFields(serviceUnitConfig *model.ServiceUnitConfig, primaryAdapterConfig *model.PrimaryAdapterSpec) ConfigValidationError {
	var adapterFieldErrors []InvalidAdapterFieldValueError
	if primaryAdapterConfig.ServerConfig != nil {
		adapterFieldErrors = ValidateServerFields(
			primaryAdapterConfig.ServerConfig,
			serviceUnitConfig.Name,
		)
	}
	if primaryAdapterConfig.ConsumerConfig != nil {
		adapterFieldErrors = ValidateConsumerFields(
			primaryAdapterConfig.ConsumerConfig,
			serviceUnitConfig.Name,
		)
	}
	if primaryAdapterConfig.RepositoryConfig != nil {
		if len(primaryAdapterConfig.Steps) > 0 {
			primaryAdapterConfig.Steps = []model.Step{} // makes sure that repository service unit config have no steps defined
			logger.Logger.Warnf(
				"Steps definition found on repository config for %s. These Steps will be ignored.\n",
				serviceUnitConfig.Name,
			)
		}
		adapterFieldErrors = ValidateRepositoryFields(
			primaryAdapterConfig.RepositoryConfig,
			serviceUnitConfig.Name,
		)
	}
	var stepFieldErrors []InvalidStepFieldValueError
	for _, step := range primaryAdapterConfig.Steps {
		errs := ValidateStepFields(step, primaryAdapterConfig, serviceUnitConfig.Name)
		stepFieldErrors = append(stepFieldErrors, errs...)
	}
	return ConfigValidationError{AdapterFieldErrors: adapterFieldErrors, StepFieldErrors: stepFieldErrors}
}

// Validate the fields of the secondary adapter configuration
func validateSecondaryAdapterConfigFields(secondaryAdapterConfig model.SecondaryAdapterConfig) ConfigValidationError {
	var errs []InvalidAdapterFieldValueError
	if secondaryAdapterConfig.InvocationConfig != nil {
		errs = ValidateInvocationFields(*secondaryAdapterConfig.InvocationConfig)
	}
	if secondaryAdapterConfig.ProducerConfig != nil {
		errs = ValidateProducerFields(*secondaryAdapterConfig.ProducerConfig)
	}
	if secondaryAdapterConfig.RepositoryConfig != nil {
		errs = ValidateRepositoryClientFields(*secondaryAdapterConfig.RepositoryConfig)
	}
	return ConfigValidationError{AdapterFieldErrors: errs}
}
