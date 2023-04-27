package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Inteface to
type Adapter interface {
	Validate() []InvalidFieldValueError
	GetId() string
}

func (sac StatelessAdapterConfig) Validate() []InvalidFieldValueError {
	validate := validator.New()
	var errs []InvalidFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidFieldValueErrors(err, sac)...)
	}

	return errs
}

func (sac StatelessAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		sac.Service,
		sac.Variant,
		sac.Action,
		sac.Route,
	)
}

func (sac StatefulAdapterConfig) Validate() []InvalidFieldValueError {
	validate := validator.New()
	var errs []InvalidFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidFieldValueErrors(err, sac)...)
	}

	return errs
}

func (sac StatefulAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s.%s",
		sac.Variant,
		sac.Action,
		sac.Size,
	)
}

func (bac BrokerAdapterConfig) Validate() []InvalidFieldValueError {
	validate := validator.New()
	var errs []InvalidFieldValueError
	err := validate.Struct(bac)
	if err != nil {
		errs = append(errs, mapInvalidFieldValueErrors(err, bac)...)
	}

	return errs
}

func (bac BrokerAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s",
		bac.Variant,
		bac.Topic,
	)
}

func (iac InternalAdapterConfig) Validate() []InvalidFieldValueError {
	validate := validator.New()
	var errs []InvalidFieldValueError
	err := validate.Struct(iac)
	if err != nil {
		errs = append(errs, mapInvalidFieldValueErrors(err, iac)...)
	}

	return errs
}

func (iac InternalAdapterConfig) GetId() string {
	return iac.Name 
}

func mapInvalidFieldValueErrors(err error, adapter Adapter) []InvalidFieldValueError {
	var errs []InvalidFieldValueError
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errs = append(errs, NewInvalidFieldValueError(fieldError.StructField(), adapter))
		}
	}
	return errs

}

func validateAdapter[T Adapter](adapter T) []InvalidFieldValueError {
	return adapter.Validate()
}
