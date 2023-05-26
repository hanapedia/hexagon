package config

import (
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/loader"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/validation"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/yaml"
)

func newConfigLoader(format string) loader.ConfigLoader {
	var configLoader loader.ConfigLoader
	switch format {
	case "yaml":
		configLoader = yaml.YamlConfigLoader{Path: "./config/service-unit.yaml"}
	default:
		configLoader = yaml.YamlConfigLoader{Path: "./config/service-unit.yaml"}
	}
	return configLoader
}

func GetConfig(format string) model.ServiceUnitConfig {

	configLoader := newConfigLoader(format)
	config, err := configLoader.Load()
	if err != nil {
		logger.Logger.Fatalf("Error loading config: %v", err)
	}

	errs := validation.ValidateServiceUnitConfigFields(&config)
	if errs.Exist() {
		errs.Print()
		logger.Logger.Fatalln("Validation failed. Aborted.")
	}
	return config
}
