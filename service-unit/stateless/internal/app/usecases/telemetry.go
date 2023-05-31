package usecases

import (
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/telemetry/tracing"
)

func InitTelemetry(name string) {
	if !config.GetEnvs().TRACING {
		return
	}
	collectorUrl := config.GetOtelCollectorUrl()
	_ = tracing.InitTracer(name, collectorUrl)
	return
}
