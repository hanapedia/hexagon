package model

import "github.com/go-playground/validator/v10"

func (suc ServiceUnitConfig) Validate() []InvalidServiceUnitFieldValueError {
	validate := validator.New()
	var errs []InvalidServiceUnitFieldValueError
	err := validate.Struct(suc)
	if err != nil {
		errs = append(errs, mapInvalidServiceUnitFieldValueErrors(err, suc)...)
	}

	return errs
}
