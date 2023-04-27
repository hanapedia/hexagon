package usecases

import (
	"log"

	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/config/yaml"
)

func newConfigLoader(format string) model.ConfigLoader {
	var configLoader model.ConfigLoader
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
		log.Fatalf("Error loading config: %v", err)
	}

	errs := model.ValidateServiceUnitConfigFields(config)
	if len(errs) > 0 {
		for _, fe := range errs {
			log.Println(fe.Error())
		}
		log.Fatalln("Validation failed. Aborted.")
	}
	return config
}
