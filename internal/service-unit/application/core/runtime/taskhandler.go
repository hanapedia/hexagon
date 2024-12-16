package runtime

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
)

func TaskSetHandler(ctx context.Context, handler *domain.PrimaryAdapterHandler) domain.TaskSetResult {
	resultCh := make(chan *secondary.SecondaryPortCallResult, len(handler.TaskSet))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for _, taskHandler := range handler.TaskSet {
			if ctx.Err() != nil {
				// stop further processing if context has been cancelled
				resultCh <- &secondary.SecondaryPortCallResult{
					Payload: nil,
					Error:   ctx.Err(),
				}
				continue
			}
			taskHandler(ctx, resultCh)
		}
	}()

	var results []*secondary.SecondaryPortCallResult
	shouldFail := false
	errCount := 0
	for i := 0; i < len(handler.TaskSet); i++ {
		result := <-resultCh
		results = append(results, result)
		if result.Error != nil {
			shouldFail = shouldFail || result.GetIsCritical()
			errCount++
		}
		// shouldFail when all adapters return err even if none of them have IsCritical set to True
		shouldFail = shouldFail || errCount == len(handler.TaskSet)
		if shouldFail {
			cancel() // cancel context immediately to prevent further processing
		}
	}
	close(resultCh)

	return domain.TaskSetResult{ShouldFail: shouldFail, TaskResults: results}
}
