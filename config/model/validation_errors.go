package model

import "fmt"

type InvalidFieldValueError struct {
	message string
}

func (e *InvalidFieldValueError) Error() string {
	return e.message
}

func NewInvalidFieldValueError(key string, adapter Adapter) InvalidFieldValueError {
	return InvalidFieldValueError{message: fmt.Sprintf("Invalid value in: %v found for key: %s", adapter, key)}
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
    FieldErrors []InvalidFieldValueError
    MappingErrors []InvalidAdapterMappingError
}
