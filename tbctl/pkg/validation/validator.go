package validation

import (
	"log"

	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/tbctl/pkg/loader"
)

func ValidateFile(path string) model.ConfigValidationError {
	serviceUnitConfig := loader.GetConfig(path)
	errs := model.ValidateServiceUnitConfigFields(serviceUnitConfig)
	if errs.Exist() {
		errs.Print()
	} else {
		log.Print("No validation error found.")
	}
	return errs
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
	if errs.Exist() {
		errs.Print()
	} else {
		log.Printf("No validation error found. Validated %v service configs.", len(serviceUnitConfigs))
	}
	return errs
}
