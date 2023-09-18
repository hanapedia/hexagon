package validation

import (
	"github.com/go-playground/validator/v10"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

func ValidateInvocationFields(sac model.InvocationConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId())...)
	}

	return errs
}

func ValidateRepositoryClientFields(sac model.RepositoryClientConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(sac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, sac.GetId())...)
	}

	return errs
}

func ValidateProducerFields(bac model.ProducerConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(bac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, bac.GetId())...)
	}

	return errs
}
