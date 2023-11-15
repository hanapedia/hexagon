package initialization

import (
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/sirupsen/logrus"
)

func InitTracing(name string) {
	if !config.GetEnvs().TRACING {
		logger.Logger.Info("Tracing is disabled.")
		return
	}
	collectorUrl := config.GetOtelCollectorUrl()
	_ = tracing.InitTracer(name, collectorUrl)
	logger.Logger.Info("Successfully initialized tracing.")
	return
}

func InitLogging() {
	if level, err := logrus.ParseLevel(config.GetEnvs().LOG_LEVEL); err == nil {
		logger.Logger.SetLevel(level)
	} else {
		logger.Logger.Warn("Failed to parse log level, resorting to info level.")
		logger.Logger.SetLevel(logrus.InfoLevel)
	}
}
