package model

import "github.com/go-playground/validator/v10"

// Config fields for Internal services
type InternalAdapterConfig struct {
	Name     string `yaml:"name" validate:"required"`
	Resource string `yaml:"resource" validate:"required,oneof=cpu memory disk network"`
	Duration string `yaml:"duration" validate:"required,oneof=small medium large"`
	Load     string `yaml:"load" validate:"required,oneof=small medium large"`
}

func (iac InternalAdapterConfig) Validate() []InvalidAdapterFieldValueError {
	validate := validator.New()
	var errs []InvalidAdapterFieldValueError
	err := validate.Struct(iac)
	if err != nil {
		errs = append(errs, mapInvalidAdapterFieldValueErrors(err, iac.GetId())...)
	}

	return errs
}

func (iac InternalAdapterConfig) GetId() string {
	return iac.Name
}

