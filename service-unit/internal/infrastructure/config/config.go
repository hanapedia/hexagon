package config

import (
	"os"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"gopkg.in/yaml.v3"
)

type YamlConfigLoader struct {
	Path string
}

func (ycl YamlConfigLoader) Load() (core.ServiceUnitConfig, error) {
	data, err := os.ReadFile(ycl.Path)

	var config core.ServiceUnitConfig
	err = yaml.Unmarshal(data, &config)

    return config, err
}
