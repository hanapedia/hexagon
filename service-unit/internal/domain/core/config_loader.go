package core

type ConfigLoader interface {
	Load() (ServiceUnitConfig, error)
}

type ServiceUnitConfig struct {
	Name           string          `yaml:"name"`
	HandlerConfigs []HandlerConfig `yaml:"handler"`
}

type HandlerConfig struct {
	Name     string `yaml:"name"`
	Protocol string `yaml:"protocol"`
	Action   string `yaml:"action"`
	Flow     []Step `yaml:"flow"`
}

type Step struct {
	AdapterId  string `yaml:"adapterId"`
	Concurrent bool   `yaml:"concurrent,omitempty"`
}
