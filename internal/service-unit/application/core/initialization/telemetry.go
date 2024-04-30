package initialization

import (
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

// InitTracing initiates distributed tracing if enabled
// TODO: encapsulate in to ports and adapters
func InitTracing(name string) *trace.TracerProvider {
	if !config.GetEnvs().TRACING {
		logger.Logger.Info("Tracing is disabled.")
		return nil
	}
	collectorUrl := config.GetOtelCollectorUrl()
	provider := tracing.InitTracer(name, collectorUrl)
	logger.Logger.Info("Initialized tracing.")
	return provider
}

func InitLogging() {
	if level, err := logrus.ParseLevel(config.GetEnvs().LOG_LEVEL); err == nil {
		logger.Logger.SetLevel(level)
	} else {
		logger.Logger.Warn("Failed to parse log level, resorting to info level.")
		logger.Logger.SetLevel(logrus.InfoLevel)
	}
}
