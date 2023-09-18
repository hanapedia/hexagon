package validation

import (
	"github.com/go-playground/validator/v10"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
)

func ValidateServiceUnitConfigFields(serviceUnitConfig *model.ServiceUnitConfig) ConfigValidationError {
	validate := validator.New()
	var errs []InvalidServiceUnitFieldValueError
	err := validate.Struct(serviceUnitConfig)
	if err != nil {
		errs = append(errs, mapInvalidServiceUnitFieldValueErrors(err, *serviceUnitConfig)...)
	}

	return ConfigValidationError{ServiceUnitFieldErrors: errs}
}

func ValidateStepFields(step model.Step, adapterConfig *model.PrimaryAdapterSpec, serviceName string) []InvalidStepFieldValueError {
	var stepFieldErrors []InvalidStepFieldValueError
	if step.AdapterConfig == nil {
		stepFieldErrors = append(stepFieldErrors, NewInvalidStepFieldValueError(adapterConfig.GetId(serviceName)))
	}

	return stepFieldErrors
}
