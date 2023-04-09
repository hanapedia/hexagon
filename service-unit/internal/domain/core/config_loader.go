package core

type ConfigLoader interface{
  Load() (ServiceUnitConfig, error)
}

type ServiceUnitConfig struct {
	Name                   string                  `yaml:"name"`
	ServerInterfaceConfigs []ServerInterfaceConfig `yaml:"interfaces"`
}

type ServerInterfaceConfig struct {
	Name     string `yaml:"name"`
	Protocol string `yaml:"protocol"`
	Action     string `yaml:"action"`
	Flow     []Step `yaml:"flow"`
}

type Step struct {
	InterfaceID string `yaml:"interfaceId"`
	Concurrent  bool  `yaml:"concurrent,omitempty"`
}

