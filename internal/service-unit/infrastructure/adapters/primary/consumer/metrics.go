package consumer

import (
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/metrics"
	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/utils"
)

func ObserveConsumerAdapterDuration(start time.Time, service string, config *v1.ConsumerConfig, shouldFail bool) {
	metrics.ObservePrimaryAdapterDuration(
		time.Since(start),
		domain.PrimaryAdapterDurationLabels{
			PrimaryLabels: domain.PrimaryLabels{
				ServiceName: service,
				Variant:     string(config.Variant),
				Topic:       config.Topic,
			},
			Status: domain.Status(utils.Btof64(shouldFail)),
		},
	)
}
