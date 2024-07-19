package resiliency

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

func NewTaskHandler(primaryAdapterId string, spec model.TaskSpec, adapter secondary.SecodaryPort) domain.TaskHandler {
	return func(ctx context.Context, resultChan chan<- *secondary.SecondaryPortCallResult) {
		var handler func(context.Context) secondary.SecondaryPortCallResult

		// TODO: build handler based on configuration such as precedence between retry and circuit breaker
		handler = WithCallTimeout(spec.GetCallTimeout(), adapter.Call)
		handler = WithRetry(spec.Resiliency.Retry, handler)
		handler = WithTaskTimeout(spec.GetTaskTimeout(), handler)
		handler = WithLogger(primaryAdapterId, adapter, handler)
		handler = WithCriticalError(spec.Resiliency.IsCritical, handler)

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
