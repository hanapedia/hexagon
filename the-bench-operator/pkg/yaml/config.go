package yaml

import (
	"os"

	"gopkg.in/yaml.v3"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
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
