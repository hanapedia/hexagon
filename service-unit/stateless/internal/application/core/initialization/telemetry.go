package initialization

import (
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/telemetry/tracing"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

func InitTelemetry(name string) {
	if !config.GetEnvs().TRACING {
		logger.Logger.Info("Tracing is disabled.")
		return
	}
	collectorUrl := config.GetOtelCollectorUrl()
	_ = tracing.InitTracer(name, collectorUrl)
	logger.Logger.Info("Successfully initialized tracing.")
	return
}
