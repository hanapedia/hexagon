package validation

import "fmt"

type InvalidFieldValueError struct {
	message string
}

func (e *InvalidFieldValueError) Error() string {
	return e.message
}

func NewInvalidFieldValueError(key string, value string) InvalidFieldValueError {
	return InvalidFieldValueError{message: fmt.Sprintf("Invalid value: %s found for key: %s", value, key )}
}

type InvalidEgressAdapterError struct {
	message string
}

func (e *InvalidEgressAdapterError) Error() string {
	return e.message
}

func NewInvalidEgressAdapterError(id string) InvalidEgressAdapterError {
	return InvalidEgressAdapterError{message: fmt.Sprintf("No matching ingress adapter found for egress adapter with id: %s", id)}
}
