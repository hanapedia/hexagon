package yaml

import (
	"os"

	k8syaml "sigs.k8s.io/yaml"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

type YamlConfigLoader struct {
	Path string
}

func (ycl YamlConfigLoader) Load() (model.ServiceUnitConfig, error) {
	data, err := os.ReadFile(ycl.Path)

	var config model.ServiceUnitConfig
	err = k8syaml.Unmarshal(data, &config)

    return config, err
}
