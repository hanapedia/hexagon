package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
)

func TaskSetHandler(ctx context.Context, handler *ports.PrimaryHandler) ports.TaskSetResult {
	resultCh := make(chan *ports.TaskResult, len(handler.TaskSet))

	for _, task := range handler.TaskSet {
		if task.Concurrent {
			go HandleTask(ctx, task, resultCh)
		} else {
			HandleTask(ctx, task, resultCh)
		}
	}

	var results []*ports.TaskResult
	shouldFail := false
	errCount := 0
	for i := 0; i < len(handler.TaskSet); i++ {
		result := <-resultCh
		results = append(results, result)
		handler.LogTaskResult(ctx, result)
		if result.Error != nil {
			shouldFail = shouldFail || result.Task.OnError.IsCritical
			errCount++
		}
		// shouldFail when all adapters return err even if none of them have IsCritical set to True
		shouldFail = shouldFail || errCount == len(handler.TaskSet)
	}
	close(resultCh)

	return ports.TaskSetResult{ShouldFail: shouldFail, TaskResults: results}
}

// HandleTask calls the task with possible retries
func HandleTask(ctx context.Context, task ports.Task, resultCh chan<- *ports.TaskResult) {
	// add 1 for the initial attempt
	maxAttempt := task.OnError.RetryMaxAttempt + 1
	var result ports.SecondaryPortCallResult

	for i := 0; i < maxAttempt; i++ {
		// derive new timeout for each calls because timeouts are meant to be per-call
		ctxWithTimeout, cancel := context.WithTimeout(ctx, task.GetTimeout())

		result = task.SecondaryPort.Call(ctxWithTimeout)
		cancel()

		if result.Error == nil {
			resultCh <- ports.NewTaskResult(task, result)
			return
		}

		backoff := task.OnError.GetNthBackoff(i + 1)
		timer := time.NewTimer(backoff)
		select {
		// check for the parent context expiration
		case <-ctx.Done():
			resultCh <- ports.NewTaskResult(task, ports.SecondaryPortCallResult{Payload: nil, Error: ctx.Err()})
			timer.Stop()
			return
		case <-timer.C:
			timer.Stop()
		}
	}
	resultCh <- ports.NewTaskResult(task, ports.SecondaryPortCallResult{Payload: nil, Error: fmt.Errorf("max retry attempt exceeded, lastError=%s", result.Error)})
}
