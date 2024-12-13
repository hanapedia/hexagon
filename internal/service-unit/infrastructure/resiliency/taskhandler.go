package resiliency

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/metrics"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

func NewTaskHandler(telCtx domain.TelemetryContext, spec *model.TaskSpec, adapter secondary.SecodaryPort) domain.TaskHandler {
	var handler CallWithContextAlias = WithUnWrapTaskContext(adapter.Call)
	var circuitBreaker CircuitBreaker = nil

	// Skip configurations if not set
	if spec.Resiliency.AdaptiveCallTimeout.Enabled {
		handler = WithAdaptiveRTOCallTimeout(spec.Resiliency.AdaptiveCallTimeout, adapter, handler)
	} else if spec.Resiliency.CallTimeout != "" {
		handler = WithCallTimeout(spec.Resiliency.GetCallTimeout(), handler)
	}

	// Log each call
	if spec.Resiliency.LogCallError {
		handler = WithLogger(telCtx, "after call", adapter, handler)
	}

	// record Call duration
	/* handler = WithCallDurationMetrics(handler) */
	// call duration should always be recorded before retry
	// if circuit breaker is applied before retry, record call duration after circuit breaker

	if spec.Resiliency.Retry.Enabled && spec.Resiliency.CircutBreaker.Enabled {
		if spec.Resiliency.CircutBreaker.CountRetries {
			// must go through the circuit breaker first to count the retries
			handler, circuitBreaker = WithCircuitBreaker(spec.Resiliency.CircutBreaker, adapter, handler)
			handler = WithCallDurationMetrics(handler, false) // record call including the circuit breaker
			handler = WithRetry(spec.Resiliency.Retry, handler)
		} else {
			handler = WithCallDurationMetrics(handler, false) // record call without circuit breaker
			handler = WithRetry(spec.Resiliency.Retry, handler)
			handler, circuitBreaker = WithCircuitBreaker(spec.Resiliency.CircutBreaker, adapter, handler)
		}
	} else {
		if spec.Resiliency.Retry.Enabled {
			handler = WithCallDurationMetrics(handler, false)
			handler = WithRetry(spec.Resiliency.Retry, handler)
		}
		if spec.Resiliency.CircutBreaker.Enabled {
			handler, circuitBreaker = WithCircuitBreaker(spec.Resiliency.CircutBreaker, adapter, handler)
			handler = WithCallDurationMetrics(handler, false)
		}
		if !spec.Resiliency.CircutBreaker.Enabled && !spec.Resiliency.Retry.Enabled {
			handler = WithCallDurationMetrics(handler, false)
		}
	}

	if spec.Resiliency.TaskTimeout != "" {
		handler = WithTaskTimeout(spec.Resiliency.GetTaskTimeout(), handler)
	}

	// record Task duration
	handler = WithTaskDurationMetrics(handler, false)

	// Always include
	handler = WithCriticalError(spec.Resiliency.IsCritical, handler)

	if spec.Resiliency.LogTaskError {
		handler = WithLogger(telCtx, "after task", adapter, handler)
	}

	// Set Gauge metrics for the adapter
	metrics.SetGaugeMetricsFromSpecs(spec.Resiliency, telCtx)

	return func(ctx context.Context, resultChan chan<- *secondary.SecondaryPortCallResult) {
		// wrap ctx with TaskContext
		taskCtx := &TaskContext{
			circuitBreaker: circuitBreaker,
			telemetryCtx:   telCtx,
			isConcurrent:   spec.Concurrent,
		}
		// Call the handler
		if spec.Concurrent {
			go func(ctx context.Context) {
				result := handler(ctx, taskCtx)
				resultChan <- &result
			}(ctx)
		} else {
			result := handler(ctx, taskCtx)
			resultChan <- &result
		}
	}
}
