package yaml

import (
	"os"

	"github.com/hanapedia/the-bench/config/model"
	"gopkg.in/yaml.v3"
)

type YamlConfigLoader struct {
	Path string
}

func (ycl YamlConfigLoader) Load() (model.ServiceUnitConfig, error) {
	data, err := os.ReadFile(ycl.Path)

	var config model.ServiceUnitConfig
	err = yaml.Unmarshal(data, &config)

    return config, err
}
