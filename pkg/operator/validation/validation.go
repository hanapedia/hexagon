package validation

import (
	model "github.com/hanapedia/the-bench/pkg/api/v1"
)

// Validate multiple service unit configs
// serviceUnitConfigs: an array of service unit configs
// returns validation error
func ValidateServiceUnitConfigs(serviceUnitConfigs []model.ServiceUnitConfig) ConfigValidationError {
	// validate service unit and adapters fields
	var configValidationError ConfigValidationError
	configValidationError.Extend(validateFieldsForAllServices(serviceUnitConfigs))
	configValidationError.Extend(validateMapping(serviceUnitConfigs))
	return configValidationError
}

// validate fields on all of the service
func validateFieldsForAllServices(serviceUnitConfigs []model.ServiceUnitConfig) ConfigValidationError {
	var configValidationError ConfigValidationError
	for _, serviceUnitConfig := range serviceUnitConfigs {
		configValidationError.Extend(ValidateFields(&serviceUnitConfig))
	}
	return configValidationError
}
