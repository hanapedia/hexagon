package resiliency

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/metrics"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

func NewTaskHandler(telCtx domain.TelemetryContext, spec model.TaskSpec, adapter secondary.SecodaryPort) domain.TaskHandler {
	var handler CallWithContextAlias = WithUnWrapTaskContext(adapter.Call)
	var circuitBreaker CircuitBreaker = nil

	// Skip configurations if not set
	if spec.Resiliency.CallTimeout != "" {
		handler = WithCallTimeout(spec.Resiliency.GetCallTimeout(), handler)
	}

	// record Call duration
	handler = WithCallDurationMetrics(handler)

	if !spec.Resiliency.Retry.Disabled && !spec.Resiliency.CircutBreaker.Disabled {
		if spec.Resiliency.CircutBreaker.CountRetries {
			// must go through the circuit breaker first to count the retries
			handler, circuitBreaker = WithCircuitBreaker(spec.Resiliency.CircutBreaker, adapter, handler)
			handler = WithRetry(spec.Resiliency.Retry, handler)
		} else {
			handler = WithRetry(spec.Resiliency.Retry, handler)
			handler, circuitBreaker = WithCircuitBreaker(spec.Resiliency.CircutBreaker, adapter, handler)
		}
	} else {
		if !spec.Resiliency.Retry.Disabled {
			handler = WithRetry(spec.Resiliency.Retry, handler)
		}
		if !spec.Resiliency.CircutBreaker.Disabled {
			handler, circuitBreaker = WithCircuitBreaker(spec.Resiliency.CircutBreaker, adapter, handler)
		}
	}

	if spec.Resiliency.TaskTimeout != "" {
		handler = WithTaskTimeout(spec.Resiliency.GetTaskTimeout(), handler)
	}

	// record Task duration
	handler = WithTaskDurationMetrics(handler)

	// Always include
	handler = WithCriticalError(spec.Resiliency.IsCritical, handler)
	handler = WithLogger(telCtx, adapter, handler)

	// Set Gauge metrics for the adapter
	metrics.SetGaugeMetricsFromSpecs(spec.Resiliency, telCtx)

	return func(ctx context.Context, resultChan chan<- *secondary.SecondaryPortCallResult) {
		// wrap ctx with TaskContext
		taskCtx := &TaskContext{
			ctx:            ctx,
			circuitBreaker: circuitBreaker,
			telemetryCtx:   telCtx,
			isConcurrent:   spec.Concurrent,
		}
		// Call the handler
		if spec.Concurrent {
			go func(ctx context.Context) {
				result := handler(taskCtx)
				resultChan <- &result
			}(ctx)
		} else {
			result := handler(taskCtx)
			resultChan <- &result
		}
	}
}
