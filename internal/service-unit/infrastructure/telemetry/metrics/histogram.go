package metrics

import (
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/domain"
)

func ObservePrimaryAdapterDuration(duration time.Duration, labels domain.PrimaryAdapterDurationLabels) {
	metrics := GetInstance()
	metrics.PrimaryAdapterDuration.
		With(labels.AsMap()).
		Observe(float64(duration.Milliseconds()))
}

func ObserveSecondaryAdapterCallDuration(duration time.Duration, labels domain.SecondaryAdapterCallDurationLabels) {
	metrics := GetInstance()
	metrics.SecondaryAdapterCallDuration.
		With(labels.AsMap()).
		Observe(float64(duration.Milliseconds()))
}

func ObserveSecondaryAdapterTaskDuration(duration time.Duration, labels domain.SecondaryAdapterTaskDurationLabels) {
	metrics := GetInstance()
	metrics.SecondaryAdapterTaskDuration.
		With(labels.AsMap()).
		Observe(float64(duration.Milliseconds()))
}
