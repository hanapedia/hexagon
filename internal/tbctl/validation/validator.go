package validation

import (
	"github.com/hanapedia/the-bench/internal/tbctl/loader"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
	v1validator "github.com/hanapedia/the-bench/pkg/operator/validation"
)

func ValidateFile(path string) v1validator.ConfigValidationError {
	serviceUnitConfig := loader.GetConfig(path)
	errs := v1validator.ValidateServiceUnitConfigFields(&serviceUnitConfig)
	if errs.Exist() {
		errs.Print()
	} else {
		logger.Logger.Infof("No validation error found.")
	}
	return errs
}

func ValidateDirectory(path string) v1validator.ConfigValidationError {
	paths, err := loader.GetYAMLFiles(path)
	if err != nil {
		logger.Logger.Errorf("Error reading from directory %s. %s", path, err)
	}

	var serviceUnitConfigs []model.ServiceUnitConfig
	for _, path = range paths {
		serviceUnitConfigs = append(serviceUnitConfigs, loader.GetConfig(path))
	}
	errs := v1validator.ValidateServiceUnitConfigs(serviceUnitConfigs)
	if errs.Exist() {
		errs.Print()
	} else {
		logger.Logger.Infof("No validation error found. Validated %v service configs.", len(serviceUnitConfigs))
	}
	return errs
}
