package resiliency

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

func NewTaskHandler(primaryAdapterId string, spec model.TaskSpec, adapter secondary.SecodaryPort) domain.TaskHandler {
	var handler func(context.Context) secondary.SecondaryPortCallResult = adapter.Call

	// Skip configurations if not set
	if spec.Resiliency.CallTimeout != "" {
		handler = WithCallTimeout(spec.Resiliency.GetCallTimeout(), handler)
	}

	if !spec.Resiliency.Retry.Disable && !spec.Resiliency.CircutBreaker.Disable {
		if spec.Resiliency.CircutBreaker.CountRetries {
			// must go through the circuit breaker first to count the retries
			handler = WithCircuitBreaker(spec.Resiliency.CircutBreaker, adapter, handler)
			handler = WithRetry(spec.Resiliency.Retry, handler)
		} else {
			handler = WithRetry(spec.Resiliency.Retry, handler)
			handler = WithCircuitBreaker(spec.Resiliency.CircutBreaker, adapter, handler)
		}
	} else {
		if !spec.Resiliency.Retry.Disable {
			handler = WithRetry(spec.Resiliency.Retry, handler)
		}
		if !spec.Resiliency.CircutBreaker.Disable {
			handler = WithCircuitBreaker(spec.Resiliency.CircutBreaker, adapter, handler)
		}
	}

	if spec.Resiliency.TaskTimeout != "" {
		handler = WithTaskTimeout(spec.Resiliency.GetTaskTimeout(), handler)
	}

	// Always include
	handler = WithCriticalError(spec.Resiliency.IsCritical, handler)
	handler = WithLogger(primaryAdapterId, adapter, handler)

	return func(ctx context.Context, resultChan chan<- *secondary.SecondaryPortCallResult) {
		// Call the handler
		if spec.Concurrent {
			go func(ctx context.Context) {
				result := handler(ctx)
				resultChan <- &result
			}(ctx)
		} else {
			result := handler(ctx)
			resultChan <- &result
		}
	}
}
