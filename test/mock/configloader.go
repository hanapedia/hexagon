package mock

import (
	"os"
	"path/filepath"

	k8syaml "sigs.k8s.io/yaml"

	"github.com/hanapedia/hexagon/internal/config/application/ports"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/test/configfiles"
)

// ConfigLoaderMock should implement ConfigLoader from
// `github.com/hanapedia/hexagon/internal/config/application/ports`
type ConfigLoaderMock struct {
	// Path contain file for loading service unit config for unit tests
	Path string
}

func NewConfigLoader(path string) ports.ConfigLoader {
	return ConfigLoaderMock{Path: path}
}

func (clm ConfigLoaderMock) Load() (*model.ServiceUnitConfig, error) {
	data, err := os.ReadFile(filepath.Join(configfiles.GetConfigFilesDir(), clm.Path))
	if err != nil {
		return nil, err
	}
	var config model.ServiceUnitConfig
	err = k8syaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
