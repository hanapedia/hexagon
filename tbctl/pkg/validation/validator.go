package validation

import (
	"log"

	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/tbctl/pkg/loader"
)

func ValidateFile(path string) ([]model.InvalidServiceUnitFieldValueError, []model.InvalidAdapterFieldValueError) {
	serviceUnitConfig := loader.GetConfig(path)
	sufe, afe := model.ValidateServiceUnitConfigFields(serviceUnitConfig)
	if len(sufe) == 0 && len(afe) == 0 {
		log.Printf("No validation error found. Validated %s", serviceUnitConfig.Name)
		return nil, nil
	}
	for _, fe := range sufe {
		log.Println(fe.Error())
	}
	for _, fe := range afe {
		log.Println(fe.Error())
	}
    return sufe, afe
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
	if len(errs.ServiceUnitFieldErrors) > 0 {
		for _, fe := range errs.ServiceUnitFieldErrors {
			log.Println(fe.Error())
		}
	}
	if len(errs.AdapterFieldErrors) > 0 {
		for _, fe := range errs.AdapterFieldErrors {
			log.Println(fe.Error())
		}
	}
	if len(errs.MappingErrors) > 0 {
		for _, me := range errs.MappingErrors {
			log.Println(me.Error())
		}
	}
	if len(errs.ServiceUnitFieldErrors) == 0 && len(errs.AdapterFieldErrors) == 0 && len(errs.MappingErrors) == 0 {
		log.Printf("No validation error found. Validated %v service configs.", len(serviceUnitConfigs))
	}
	return errs
}
