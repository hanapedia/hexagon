package usecases

import (
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
)

func GetConfig(format string) model.ServiceUnitConfig {
	serviceUnitConfig := config.GetConfig(format)
	return serviceUnitConfig
}
