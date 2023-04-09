package config

import (
	"io/ioutil"
	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"gopkg.in/yaml.v3"
)

type YamlConfigLoader struct {
	path string
}

func (ycl YamlConfigLoader) Load() (core.ServiceUnitConfig, error) {
	data, err := ioutil.ReadFile(ycl.path)

	var config core.ServiceUnitConfig
	err = yaml.Unmarshal(data, &config)

    return config, err
}
