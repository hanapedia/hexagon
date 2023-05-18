package validation

import (
	"github.com/go-playground/validator/v10"

	thebenchv1 "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

func ValidateStatelessIngressAdapterFields(sac thebenchv1.StatelessIngressAdapterConfig, serviceName string) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId(serviceName))...)
	}

	return errs
}

func ValidateStatefulIngressAdapterFields(sac thebenchv1.StatefulIngressAdapterConfig, serviceName string) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId(serviceName))...)
	}

	return errs
}

func ValidateBrokerIngressAdapterFields(bac thebenchv1.BrokerIngressAdapterConfig, serviceName string) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(bac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, bac.GetId(serviceName))...)
	}

	return errs
}
