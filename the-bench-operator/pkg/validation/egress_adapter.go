package validation

import (
	"github.com/go-playground/validator/v10"

	thebenchv1 "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

func ValidateStatelessEgressAdapterFields(sac thebenchv1.StatelessEgressAdapterConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId())...)
	}

	return errs
}

func ValidateStatefulEgressAdapterFields(sac thebenchv1.StatefulEgressAdapterConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId())...)
	}

	return errs
}

func ValidateBrokerEgressAdapterFields(bac thebenchv1.BrokerEgressAdapterConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(bac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, bac.GetId())...)
	}

	return errs
}
