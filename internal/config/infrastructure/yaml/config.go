package yaml

import (
	"os"

	k8syaml "sigs.k8s.io/yaml"

	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

type YamlConfigLoader struct {
	Path string
}

func (ycl YamlConfigLoader) Load() (*model.ServiceUnitConfig, error) {
	data, err := os.ReadFile(ycl.Path)
	if err != nil {
		logger.Logger.
			WithField("path", ycl.Path).
			WithField("err", err).
			Error("Failed to read config file ")
		return nil, err
	}

	var config model.ServiceUnitConfig
	err = k8syaml.UnmarshalStrict(data, &config)
	if err != nil {
		logger.Logger.
			WithField("path", ycl.Path).
			WithField("err", err).
			Error("Failed to load config from yaml ")
		return nil, err
	}

	if config.Kind == model.ClusterConfigKind {
		logger.Logger.
			WithField("path", ycl.Path).
			WithField("err", "Cannot load ClusterConfig as ServiceUnit.").
			Error("Failed to load config from yaml ")
		return nil, err
	}

	if config.Kind == model.ServiceUnitKind {
		logger.Logger.
			WithField("path", ycl.Path).
			Warn("Attempting to parse a yaml file without valid `kind` field.")
	}

	return &config, err
}

func (ycl YamlConfigLoader) LoadClusterConfig() (*model.ClusterConfig, error) {
	data, err := os.ReadFile(ycl.Path)
	if os.IsNotExist(err) {
		logger.Logger.
			WithField("path", ycl.Path).
			Infof("%s not found. Skipping.", constants.CLUSTER_CONFIG_FILE_NAME)
		clusterConfig := model.NewClusterConfig()
		return &clusterConfig, nil
	}
	if err != nil {
		logger.Logger.
			WithField("path", ycl.Path).
			WithField("err", err).
			Error("Failed to read config file ")
		return nil, err
	}

	var config model.ClusterConfig = model.NewClusterConfig()
	err = k8syaml.Unmarshal(data, &config)
	if err != nil {
		logger.Logger.
			WithField("path", ycl.Path).
			WithField("err", err).
			Error("Failed to load config from yaml ")
		return nil, err
	}

	if config.Kind != model.ClusterConfigKind {
		logger.Logger.
			WithField("path", ycl.Path).
			WithField("err", "Invalid `kind` field.").
			Error("Failed to load cluster config from yaml ")
		return nil, err
	}

	return &config, err
}
