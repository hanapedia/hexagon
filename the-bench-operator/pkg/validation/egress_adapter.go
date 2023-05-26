package validation

import (
	"github.com/go-playground/validator/v10"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

func ValidateStatelessEgressAdapterFields(sac model.StatelessEgressAdapterConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId())...)
	}

	return errs
}

func ValidateStatefulEgressAdapterFields(sac model.StatefulEgressAdapterConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId())...)
	}

	return errs
}

func ValidateBrokerEgressAdapterFields(bac model.BrokerEgressAdapterConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(bac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, bac.GetId())...)
	}

	return errs
}
