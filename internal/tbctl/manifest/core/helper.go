package core

import (
	"fmt"
	"os"

	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

func FormatManifest(manifest []byte) string {
	return fmt.Sprintf("---\n%s\n", manifest)
}

// CreateFile create and open file in append
func CreateFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
}

func GetFilePath(dir, name, kind string) string {
	return fmt.Sprintf("%s/%s-%s.yaml", dir, name, kind)
}

func HasRepositoryAdapter(config *model.ServiceUnitConfig) bool {
	for _, primaryAdapterConfig := range config.AdapterConfigs {
		if primaryAdapterConfig.RepositoryConfig != nil {
			return true
		}
	}
	return false
}

func GetRepositoryAdapter(config *model.ServiceUnitConfig) *model.RepositoryConfig {
	for _, primaryAdapterConfig := range config.AdapterConfigs {
		if primaryAdapterConfig.RepositoryConfig != nil {
			return primaryAdapterConfig.RepositoryConfig
		}
	}
	return nil
}

func HasConsumerAdapters(config *model.ServiceUnitConfig) bool {
	for _, primaryAdapterConfig := range config.AdapterConfigs {
		if primaryAdapterConfig.ConsumerConfig != nil {
			return true
		}
	}
	return false
}

func GetConsumerAdapters(config *model.ServiceUnitConfig) []model.ConsumerConfig {
	var configs []model.ConsumerConfig
	for _, primaryAdapterConfig := range config.AdapterConfigs {
		if primaryAdapterConfig.ConsumerConfig != nil {
			configs = append(configs, *primaryAdapterConfig.ConsumerConfig)
		}
	}
	return configs
}

func HasGatewayConfig(config *model.ServiceUnitConfig) bool {
	return config.DeploymentSpec.Gateway != nil
}

