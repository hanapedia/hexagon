package usecases

import (
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/config/yaml"
)

func NewConfigLoader(format string) model.ConfigLoader {
	var configLoader model.ConfigLoader
	switch format {
	case "yaml":
		configLoader = yaml.YamlConfigLoader{Path: "./config/service-unit.yaml"}
	default:
		configLoader = yaml.YamlConfigLoader{Path: "./config/service-unit.yaml"}
	}
	return configLoader
}
