package config

import (
	"github.com/hanapedia/the-bench/config/loader"
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/config/validation"
	"github.com/hanapedia/the-bench/config/logger"
	"github.com/hanapedia/the-bench/config/yaml"
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
		logger.PrintErrors(errs)
		logger.Logger.Fatalln("Validation failed. Aborted.")
	}
	return config
}

