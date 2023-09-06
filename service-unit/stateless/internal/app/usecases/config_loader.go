package usecases

import (
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
)

// GetConfig parses service-unit config via given formant
func GetConfig(format string) model.ServiceUnitConfig {
	serviceUnitConfig := config.GetConfig(format)
	return serviceUnitConfig
}
