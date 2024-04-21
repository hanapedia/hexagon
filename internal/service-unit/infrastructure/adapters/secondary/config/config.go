package config

import (
	"github.com/hanapedia/hexagon/internal/config/application/ports"
	"github.com/hanapedia/hexagon/internal/config/infrastructure/yaml"
)

func NewConfigLoader(format string) ports.ConfigLoader {
	var configLoader ports.ConfigLoader
	switch format {
	case "yaml":
		configLoader = yaml.YamlConfigLoader{Path: "./config/service-unit.yaml"}
	default:
		configLoader = yaml.YamlConfigLoader{Path: "./config/service-unit.yaml"}
	}
	return configLoader
}

