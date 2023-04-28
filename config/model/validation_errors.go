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

func NewInvalidAdapterFieldValueError(key string, adapter Adapter, message string) InvalidAdapterFieldValueError {
	return InvalidAdapterFieldValueError{message: fmt.Sprintf("Invalid value in adapter: %v for key: %s. %s", adapter.GetId(), key, message)}
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

type ConfigValidationError struct {
	ServiceUnitFieldErrors []InvalidServiceUnitFieldValueError
	AdapterFieldErrors     []InvalidAdapterFieldValueError
	MappingErrors          []InvalidAdapterMappingError
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

func mapInvalidAdapterFieldValueErrors(err error, adapter Adapter) []InvalidAdapterFieldValueError {
	var errs []InvalidAdapterFieldValueError
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errs = append(errs, NewInvalidAdapterFieldValueError(fieldError.StructField(), adapter, fieldError.Error()))
		}
	}
	return errs
}
