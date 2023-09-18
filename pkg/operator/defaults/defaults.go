package defaults

import model "github.com/hanapedia/the-bench/pkg/api/v1"

func SetDefaults(serviceUnitConfig *model.ServiceUnitConfig) {
	if serviceUnitConfig.Version == "" {
		serviceUnitConfig.Version = "latest"
	}
}
