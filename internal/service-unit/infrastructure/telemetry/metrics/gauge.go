package metrics

import (
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/utils"
)

func SetPrimaryAdapterInProgress(op domain.GaugeOp, labels domain.PrimaryAdapterInProgressLabels) {
	metrics := GetInstance()
	switch op {
	case domain.INC:
		metrics.PrimaryAdapterInProgress.With(labels.AsMap()).Inc()
	case domain.DEC:
		metrics.PrimaryAdapterInProgress.With(labels.AsMap()).Dec()
	}
}

// SetGaugeMetricsFromSpecs sets gauge metrics from resiliency spec.
// `With` can panic, however, it shouldn't if the labels match what's defined in Metrics instance.
func SetGaugeMetricsFromSpecs(spec model.ResiliencySpec, telemetryCtx domain.TelemetryContext) {
	metrics := GetInstance()

	timeoutGaugeLabels := domain.TimeoutGaugeLabels{Ctx: telemetryCtx}.AsMap()
	circuitBreakerGaugeLabels := domain.CircuitBreakerGaugeLabels{Ctx: telemetryCtx}.AsMap()
	retryGaugeLabels := domain.RetryGaugeLabels{Ctx: telemetryCtx, BackoffPolicy: spec.Retry.BackoffPolicy}.AsMap()

	metrics.CallTimeout.
		With(timeoutGaugeLabels).
		Set(float64(spec.GetCallTimeout().Milliseconds()))

	metrics.TaskTimeout.
		With(timeoutGaugeLabels).
		Set(float64(spec.GetTaskTimeout().Milliseconds()))

	metrics.CircuitBreakerDisabled.
		With(circuitBreakerGaugeLabels).
		Set(utils.Btof64(spec.CircutBreaker.Disabled))

	metrics.CircuitBreakerCountRetries.
		With(circuitBreakerGaugeLabels).
		Set(utils.Btof64(spec.CircutBreaker.CountRetries))

	metrics.CircuitBreakerIntervalSecs.
		With(circuitBreakerGaugeLabels).
		Set(spec.CircutBreaker.GetInterval().Seconds())

	metrics.CircuitBreakerMaxRequests.
		With(circuitBreakerGaugeLabels).
		Set(float64(spec.CircutBreaker.MaxRequests))

	metrics.CircuitBreakerMinRequests.
		With(circuitBreakerGaugeLabels).
		Set(float64(spec.CircutBreaker.MinRequests))

	metrics.CircuitBreakerTimeout.
		With(circuitBreakerGaugeLabels).
		Set(spec.CircutBreaker.GetTimeout().Seconds())

	metrics.CircuitBreakerRatio.
		With(circuitBreakerGaugeLabels).
		Set(spec.CircutBreaker.Ratio)

	metrics.CircuitBreakerConsecutiveFails.
		With(circuitBreakerGaugeLabels).
		Set(float64(spec.CircutBreaker.ConsecutiveFails))

	metrics.RetryDisabled.
		With(retryGaugeLabels).
		Set(utils.Btof64(spec.Retry.Disabled))

	metrics.RetryMaxAttempt.
		With(retryGaugeLabels).
		Set(float64(spec.Retry.MaxAttempt))

	metrics.RetryInitialBackoff.
		With(retryGaugeLabels).
		Set(float64(spec.Retry.GetInitialBackoff().Milliseconds()))
}
