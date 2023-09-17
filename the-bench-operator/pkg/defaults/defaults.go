package defaults

import model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"

func SetDefaults(serviceUnitConfig *model.ServiceUnitConfig) {
	if serviceUnitConfig.Version == "" {
		serviceUnitConfig.Version = "latest"
	}
}
