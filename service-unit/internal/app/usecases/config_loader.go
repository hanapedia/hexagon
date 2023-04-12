package usecases

import (
	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/config"
)

func NewConfigLoader(format string) core.ConfigLoader {
	var configLoader core.ConfigLoader
	switch format {
	case "yaml":
		configLoader = config.YamlConfigLoader{Path: "./config/service-unit.yaml"}
	default:
		configLoader = config.YamlConfigLoader{Path: "./config/service-unit.yaml"}
	}
	return configLoader
}
