package server

import (
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/metrics"
	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/utils"
)

func ObserveServerAdapterDuration(duration time.Duration, service string, config *v1.ServerConfig, shouldFail bool) {
	metrics.ObservePrimaryAdapterDuration(
		duration,
		domain.PrimaryAdapterDurationLabels{
			PrimaryLabels: domain.PrimaryLabels{
				ServiceName: service,
				Variant:     string(config.Variant),
				Route:       config.Route,
				Action:      string(config.Action),
			},
			Status: domain.Status(utils.Btof64(shouldFail)),
		},
	)
}
