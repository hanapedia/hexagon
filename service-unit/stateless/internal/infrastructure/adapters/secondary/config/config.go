package config

import (
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/loader"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/yaml"
)

func NewConfigLoader(format string) loader.ConfigLoader {
	var configLoader loader.ConfigLoader
	switch format {
	case "yaml":
		configLoader = yaml.YamlConfigLoader{Path: "./config/service-unit.yaml"}
	default:
		configLoader = yaml.YamlConfigLoader{Path: "./config/service-unit.yaml"}
	}
	return configLoader
}

