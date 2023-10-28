package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

// An object to hold all the errors of different types
type ConfigValidationError struct {
	ServiceUnitFieldErrors []InvalidServiceUnitFieldValueError
	AdapterFieldErrors     []InvalidAdapterFieldValueError
	MappingErrors          []InvalidAdapterMappingError
	StepFieldErrors        []InvalidStepFieldValueError
}

type InvalidServiceUnitFieldValueError struct {
	message string
}

type InvalidAdapterFieldValueError struct {
	message string
}

type InvalidAdapterMappingError struct {
	message string
}

type InvalidStepFieldValueError struct {
	message string
}

func (cve ConfigValidationError) Exist() bool {
	return (len(cve.ServiceUnitFieldErrors) > 0 ||
		len(cve.AdapterFieldErrors) > 0 ||
		len(cve.MappingErrors) > 0 ||
		len(cve.StepFieldErrors) > 0)
}

func (cve *ConfigValidationError) Extend(other ConfigValidationError) {
	cve.ServiceUnitFieldErrors = append(cve.ServiceUnitFieldErrors, other.ServiceUnitFieldErrors...)
	cve.AdapterFieldErrors = append(cve.AdapterFieldErrors, other.AdapterFieldErrors...)
	cve.MappingErrors = append(cve.MappingErrors, other.MappingErrors...)
	cve.StepFieldErrors = append(cve.StepFieldErrors, other.StepFieldErrors...)
}

func (cve ConfigValidationError) Print() {
	for _, err := range cve.ServiceUnitFieldErrors {
		logger.Logger.Errorf(err.Error())
	}
	for _, err := range cve.AdapterFieldErrors {
		logger.Logger.Errorf(err.Error())
	}
	for _, err := range cve.MappingErrors {
		logger.Logger.Errorf(err.Error())
	}
	for _, err := range cve.StepFieldErrors {
		logger.Logger.Errorf(err.Error())
	}
}

func (e *InvalidServiceUnitFieldValueError) Error() string {
	return e.message
}

func (e *InvalidAdapterFieldValueError) Error() string {
	return e.message
}

func (e *InvalidAdapterMappingError) Error() string {
	return e.message
}

func (e *InvalidStepFieldValueError) Error() string {
	return e.message
}

func NewInvalidServiceUnitFieldValueError(key string, serviceUnitConfig model.ServiceUnitConfig, message string) InvalidServiceUnitFieldValueError {
	return InvalidServiceUnitFieldValueError{message: fmt.Sprintf("Invalid value in service unit definition: %v for key: %s. %s", serviceUnitConfig.Name, key, message)}
}

func NewInvalidAdapterFieldValueError(key string, adapterId string, message string) InvalidAdapterFieldValueError {
	return InvalidAdapterFieldValueError{message: fmt.Sprintf("Invalid value in adapter: %v for key: %s. %s", adapterId, key, message)}
}

func NewInvalidSecondaryAdapterError(id string) InvalidAdapterMappingError {
	return InvalidAdapterMappingError{message: fmt.Sprintf("No matching primary adapter found for secondary adapter with id: %s", id)}
}

func NewInvalidStepFieldValueError(id string) InvalidStepFieldValueError {
	return InvalidStepFieldValueError{message: fmt.Sprintf("No secondary adapter config found on one of steps on primary adapter with id: %s.", id)}
}

func mapInvalidServiceUnitFieldValueErrors(err error, serviceUnitConfig model.ServiceUnitConfig) []InvalidServiceUnitFieldValueError {
	var errs []InvalidServiceUnitFieldValueError
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errs = append(errs, NewInvalidServiceUnitFieldValueError(fieldError.StructField(), serviceUnitConfig, fieldError.Error()))
		}
	}
	return errs
}

func mapInvalidAdapterFieldValueErrors(err error, adapterId string) []InvalidAdapterFieldValueError {
	var errs []InvalidAdapterFieldValueError
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errs = append(errs, NewInvalidAdapterFieldValueError(fieldError.StructField(), adapterId, fieldError.Error()))
		}
	}
	return errs
}
