package core

import "github.com/hanapedia/the-bench/service-unit/pkg/constants"

type ConfigLoader interface {
	Load() (ServiceUnitConfig, error)
}

type ServiceUnitConfig struct {
	Name           string          `yaml:"name"`
	HandlerConfigs []HandlerConfig `yaml:"handler"`
}

type HandlerConfig struct {
	Name     string `yaml:"name"`
	Protocol constants.AdapterProtocol `yaml:"protocol"`
	Action   string `yaml:"action"`
	Steps     []Step `yaml:"flow"`
}

type Step struct {
	AdapterId  string `yaml:"adapterId"`
	Concurrent bool   `yaml:"concurrent,omitempty"`
}
