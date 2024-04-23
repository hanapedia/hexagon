package mock

import (
	k8syaml "sigs.k8s.io/yaml"

	"github.com/hanapedia/hexagon/internal/config/application/ports"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

// ConfigLoaderMock should implement ConfigLoader from
// `github.com/hanapedia/hexagon/internal/config/application/ports`
type ConfigLoaderMock struct {
	// path contain file for loading service unit config for unit tests
	data []byte
}

func NewConfigLoader(data string) ports.ConfigLoader {
	return ConfigLoaderMock{data: []byte(data)}
}

// Load for mock config loader does not validate cofiguration for separation of concerns
func (clm ConfigLoaderMock) Load() (*model.ServiceUnitConfig, error) {
	var config model.ServiceUnitConfig
	err := k8syaml.Unmarshal(clm.data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
