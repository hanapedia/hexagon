package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/hanapedia/the-bench/config/constants"
)

type EgressAdapter interface {
	Validate() []InvalidAdapterFieldValueError
	GetId() string
}

func ValidateEgressAdapter[T EgressAdapter](adapter T) []InvalidAdapterFieldValueError {
	return adapter.Validate()
}

// Config fields for stateful services
type StatelessEgressAdapterConfig struct {
	Variant constants.StatelessAdapterVariant `yaml:"variant" validate:"required,oneof=rest grpc"`
	Service string                            `yaml:"service,omitempty" validate:"required"`
	Action  constants.Action                  `yaml:"action" validate:"required,oneof=read write"`
	Route   string                            `yaml:"route" validate:"required"`
}

func (sac StatelessEgressAdapterConfig) Validate() []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId())...)
	}

	return errs
}

func (sac StatelessEgressAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		sac.Service,
		sac.Variant,
		sac.Action,
		sac.Route,
	)
}

// Config fields for stateful services
type StatefulEgressAdapterConfig struct {
	Name    string                           `yaml:"name" validate:"required"`
	Variant constants.StatefulAdapterVariant `yaml:"variant" validate:"required,oneof=mongo postgre"`
	Action  constants.Action                 `yaml:"action" validate:"omitempty,oneof=read write"`
	Size    string                           `yaml:"size" validate:"omitempty,oneof=small medium large"`
}

func (sac StatefulEgressAdapterConfig) Validate() []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId())...)
	}

	return errs
}

func (sac StatefulEgressAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s",
		sac.Variant,
		sac.Name,
	)
}

// Config fields for Brokers
type BrokerEgressAdapterConfig struct {
	Variant constants.BrokerAdapterVariant `yaml:"variant" validate:"required,oneof=kafka rabbitmq pulsar"`
	Topic   string                         `yaml:"topic" validate:"required"`
}

func (bac BrokerEgressAdapterConfig) Validate() []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(bac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, bac.GetId())...)
	}

	return errs
}

func (bac BrokerEgressAdapterConfig) GetId() string {
	return fmt.Sprintf(
		"%s.%s",
		bac.Variant,
		bac.Topic,
	)
}
