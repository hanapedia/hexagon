package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/hanapedia/the-bench/config/constants"
)

type IngressAdapter interface {
	Validate(string) []InvalidAdapterFieldValueError
	GetId(string) string
}

func ValidateIngressAdapter[T IngressAdapter](serviceName string, adapter T) []InvalidAdapterFieldValueError {
	return adapter.Validate(serviceName)
}

// Config fields for stateful services
type StatelessIngressAdapterConfig struct {
	Variant constants.StatelessAdapterVariant `yaml:"variant" validate:"required,oneof=rest grpc"`
	Action  constants.Action                  `yaml:"action" validate:"required,oneof=read write"`
	Route   string                            `yaml:"route" validate:"required"`
}

func (sac StatelessIngressAdapterConfig) Validate(serviceName string) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId(serviceName))...)
	}

	return errs
}

func (sac StatelessIngressAdapterConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s.%s.%s",
		serviceName,
		sac.Variant,
		sac.Action,
		sac.Route,
	)
}

// Config fields for stateful services
type StatefulIngressAdapterConfig struct {
	Variant constants.StatefulAdapterVariant `yaml:"variant" validate:"required,oneof=mongo postgre"`
}

func (sac StatefulIngressAdapterConfig) Validate(serviceName string) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId(serviceName))...)
	}

	return errs
}

func (sac StatefulIngressAdapterConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s",
		sac.Variant,
		serviceName,
	)
}

// Config fields for Brokers
type BrokerIngressAdapterConfig struct {
	Variant constants.BrokerAdapterVariant `yaml:"variant" validate:"required,oneof=kafka rabbitmq pulsar"`
	Topic   string                         `yaml:"topic" validate:"required"`
}

func (bac BrokerIngressAdapterConfig) Validate(serviceName string) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(bac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, bac.GetId(serviceName))...)
	}

	return errs
}

func (bac BrokerIngressAdapterConfig) GetId(serviceName string) string {
	return fmt.Sprintf(
		"%s.%s",
		bac.Variant,
		bac.Topic,
	)
}
