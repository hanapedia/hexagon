package yaml

import (
	"os"

	k8syaml "sigs.k8s.io/yaml"

	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

type YamlConfigLoader struct {
	Path string
}

func (ycl YamlConfigLoader) Load() (model.ServiceUnitConfig, error) {
	data, err := os.ReadFile(ycl.Path)
	if err != nil {
		logger.Logger.Fatal("Failed to read config file ", "path=", ycl.Path, "err=", err)
	}

	var config model.ServiceUnitConfig
	err = k8syaml.Unmarshal(data, &config)
	if err != nil {
		logger.Logger.Fatal("Failed to load config from yaml ", "path=", ycl.Path, "err=", err)
	}

    return config, err
}
