package validation

import (
	"github.com/go-playground/validator/v10"

	thebenchv1 "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

func ValidateInternalAdapterConfigFields(iac thebenchv1.InternalAdapterConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(iac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, iac.GetId())...)
	}

	return errs
}
