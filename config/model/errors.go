package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type InvalidServiceUnitFieldValueError struct {
	message string
}

func (e *InvalidServiceUnitFieldValueError) Error() string {
	return e.message
}

func NewInvalidServiceUnitFieldValueError(key string, serviceUnitConfig ServiceUnitConfig, message string) InvalidServiceUnitFieldValueError {
	return InvalidServiceUnitFieldValueError{message: fmt.Sprintf("Invalid value in service unit definition: %v for key: %s. %s", serviceUnitConfig.Name, key, message)}
}

type InvalidAdapterFieldValueError struct {
	message string
}

func (e *InvalidAdapterFieldValueError) Error() string {
	return e.message
}

func NewInvalidAdapterFieldValueError(key string, adapterId string, message string) InvalidAdapterFieldValueError {
	return InvalidAdapterFieldValueError{message: fmt.Sprintf("Invalid value in adapter: %v for key: %s. %s", adapterId, key, message)}
}

type InvalidAdapterMappingError struct {
	message string
}

func (e *InvalidAdapterMappingError) Error() string {
	return e.message
}

func NewInvalidEgressAdapterError(id string) InvalidAdapterMappingError {
	return InvalidAdapterMappingError{message: fmt.Sprintf("No matching ingress adapter found for egress adapter with id: %s", id)}
}

type InvalidStepFieldValueError struct {
	message string
}

func (e *InvalidStepFieldValueError) Error() string {
	return e.message
}

func NewInvalidStepFieldValueError(id string) InvalidStepFieldValueError {
	return InvalidStepFieldValueError{message: fmt.Sprintf("No egress adapter config found on one of steps on ingress adapter with id: %s.", id)}
}

type ConfigValidationError struct {
	ServiceUnitFieldErrors []InvalidServiceUnitFieldValueError
	AdapterFieldErrors     []InvalidAdapterFieldValueError
	MappingErrors          []InvalidAdapterMappingError
	StepFieldErrors        []InvalidStepFieldValueError
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

func mapInvalidServiceUnitFieldValueErrors(err error, serviceUnitConfig ServiceUnitConfig) []InvalidServiceUnitFieldValueError {
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
