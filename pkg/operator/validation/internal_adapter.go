package validation

import (
	"github.com/go-playground/validator/v10"

	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

func ValidateInternalAdapterConfigFields(iac model.StressorConfig) []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(iac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, iac.GetId())...)
	}

	return errs
}
