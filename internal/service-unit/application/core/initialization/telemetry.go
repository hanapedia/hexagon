package initialization

import (
	"context"
	"sync"

	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/sirupsen/logrus"
)

// InitTracing initiates distributed tracing if enabled
// TODO: encapsulate in to ports and adapters
func InitTracing(name string, ctx context.Context, wg *sync.WaitGroup) {
	if !config.GetEnvs().TRACING {
		logger.Logger.Info("Tracing is disabled.")
	}
	wg.Add(1)

	collectorUrl := config.GetOtelCollectorUrl()
	provider := tracing.InitTracer(name, collectorUrl)
	logger.Logger.Info("Initialized tracing.")

	go func() {
		<-ctx.Done()
		logger.Logger.Infof("Context cancelled. Trace provider shutting down.")
		provider.Shutdown(context.Background())
		wg.Done()
	}()
}

func InitLogging() {
	if level, err := logrus.ParseLevel(config.GetEnvs().LOG_LEVEL); err == nil {
		logger.Logger.SetLevel(level)
	} else {
		logger.Logger.Warn("Failed to parse log level, resorting to info level.")
		logger.Logger.SetLevel(logrus.InfoLevel)
	}
}
