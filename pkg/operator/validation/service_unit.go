package validation

import (
	"github.com/go-playground/validator/v10"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
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

func ValidateStepFields(task model.TaskSpec, adapterConfig *model.PrimaryAdapterSpec, serviceName string) []InvalidStepFieldValueError {
	var taskFieldErrors []InvalidStepFieldValueError
	if task.AdapterConfig == nil {
		taskFieldErrors = append(taskFieldErrors, NewInvalidStepFieldValueError(adapterConfig.GetId(serviceName)))
	}

	return taskFieldErrors
}
