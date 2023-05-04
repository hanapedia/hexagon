package validation

import (
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/config/logger"
	validator "github.com/hanapedia/the-bench/config/validation"
	"github.com/hanapedia/the-bench/tbctl/pkg/loader"
)

func ValidateFile(path string) model.ConfigValidationError {
	serviceUnitConfig := loader.GetConfig(path)
	errs := validator.ValidateServiceUnitConfigFields(&serviceUnitConfig)
	if errs.Exist() {
		logger.PrintErrors(errs)
	} else {
		logger.Logger.Infof("No validation error found.")
	}
	return errs
}

func ValidateDirectory(path string) model.ConfigValidationError {
	paths, err := loader.GetYAMLFiles(path)
	if err != nil {
		logger.Logger.Errorf("Error reading from directory %s. %s", path, err)
	}

	var serviceUnitConfigs []model.ServiceUnitConfig
	for _, path = range paths {
		serviceUnitConfigs = append(serviceUnitConfigs, loader.GetConfig(path))
	}
	errs := validator.ValidateServiceUnitConfigs(&serviceUnitConfigs)
	if errs.Exist() {
		logger.PrintErrors(errs)
	} else {
		logger.Logger.Infof("No validation error found. Validated %v service configs.", len(serviceUnitConfigs))
	}
	return errs
}
