package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Inteface to
type Adapter interface {
	Validate() []InvalidAdapterFieldValueError
	GetId() string
}

func (sac StatelessAdapterConfig) Validate() []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac)...)
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

func (sac StatefulAdapterConfig) Validate() []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac)...)
	}

	return errs
}

func (sac StatefulAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s.%s",
		sac.Variant,
		sac.Name,
		sac.Size,
	)
}

func (bac BrokerAdapterConfig) Validate() []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(bac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, bac)...)
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

func (iac InternalAdapterConfig) Validate() []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(iac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, iac)...)
	}

	return errs
}

func (iac InternalAdapterConfig) GetId() string {
	return iac.Name 
}

func validateAdapter[T Adapter](adapter T) []InvalidAdapterFieldValueError {
	return adapter.Validate()
}
