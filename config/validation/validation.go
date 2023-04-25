package validation

import (
	"github.com/hanapedia/the-bench/config/model"
)

func ValidateConfig(serviceUnitConfigs []model.ServiceUnitConfig) []error {
	baseLen := len(serviceUnitConfigs)
	validationErrors := make([]error, baseLen)
	serviceAdapterMap := make(map[string][]string)

	for _, serviceUnitConfig := range serviceUnitConfigs {
		serviceAdapterMap[serviceUnitConfig.Name] = make([]string, baseLen)
		for _, handlerConfig := range serviceUnitConfig.HandlerConfigs {

		}
	}
}
