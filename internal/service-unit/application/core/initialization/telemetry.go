package initialization

import (
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
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
