package validation

import (
	"log"

	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/tbctl/pkg/loader"
)

func ValidateFile(path string) []model.InvalidFieldValueError {
	serviceUnitConfig := loader.GetConfig(path)
	errs := model.ValidateServiceUnitConfigFields(serviceUnitConfig)
	if len(errs) > 0 {
		for _, fe := range errs {
			log.Println(fe.Error())
		}
		return errs
	}
	log.Printf("No validation error found. Validated %s", serviceUnitConfig.Name)
	return nil
}

func ValidateDirectory(path string) model.ConfigValidationError {
	paths, err := loader.GetYAMLFiles(path)
	if err != nil {
		log.Fatalf("Error reading from directory %s. %s", path, err)
	}

	var serviceUnitConfigs []model.ServiceUnitConfig
	for _, path = range paths {
		serviceUnitConfigs = append(serviceUnitConfigs, loader.GetConfig(path))
	}
	errs := model.ValidateServiceUnitConfigs(serviceUnitConfigs)
	if len(errs.FieldErrors) > 0 {
		for _, fe := range errs.FieldErrors {
			log.Println(fe.Error())
		}
	}
	if len(errs.MappingErrors) > 0 {
		for _, me := range errs.MappingErrors {
			log.Println(me.Error())
		}
	}
	if len(errs.FieldErrors) == 0 && len(errs.MappingErrors) == 0 {
		log.Printf("No validation error found. Validated %v service configs.", len(serviceUnitConfigs))
	}
    return errs
}
