package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/primary"
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
)

func TaskSetHandler(ctx context.Context, handler *primary.PrimaryHandler) primary.TaskSetResult {
	resultCh := make(chan *primary.TaskResult, len(handler.TaskSet))

	for _, task := range handler.TaskSet {
		// Add timeout to context
		taskCtx, taskCancel := context.WithTimeout(ctx, task.GetTaskTimeout())
		defer taskCancel()

		if task.Concurrent {
			go HandleTask(taskCtx, task, resultCh)
		} else {
			HandleTask(taskCtx, task, resultCh)
		}
	}

	var results []*primary.TaskResult
	shouldFail := false
	errCount := 0
	for i := 0; i < len(handler.TaskSet); i++ {
		result := <-resultCh
		results = append(results, result)
		primary.LogTaskResult(ctx, handler.GetId(), result) //TODO: reimport
		if result.Error != nil {
			shouldFail = shouldFail || result.Task.OnError.IsCritical
			errCount++
		}
		// shouldFail when all adapters return err even if none of them have IsCritical set to True
		shouldFail = shouldFail || errCount == len(handler.TaskSet)
	}
	close(resultCh)

	return primary.TaskSetResult{ShouldFail: shouldFail, TaskResults: results}
}

// HandleTask calls the task with possible retries
func HandleTask(taskCtx context.Context, task primary.Task, resultCh chan<- *primary.TaskResult) {
	// add 1 for the initial attempt
	maxAttempt := task.OnError.Retry.MaxAttempt + 1
	var result secondary.SecondaryPortCallResult

	for i := 0; i < maxAttempt; i++ {
		if i > 0 {
			backoff := task.OnError.Retry.GetNthBackoff(i)
			timer := time.NewTimer(backoff)
			select {
			// check for the parent context expiration
			case <-taskCtx.Done():
				resultCh <- primary.NewTaskResult(task, secondary.SecondaryPortCallResult{Payload: nil, Error: taskCtx.Err()})
				timer.Stop()
				return
			case <-timer.C:
				timer.Stop()
			}
		}

		// derive new timeout for each calls
		callCtx, callCancel := context.WithTimeout(taskCtx, task.GetCallTimeout())

		result = task.SecondaryPort.Call(callCtx)
		callCancel()

		if result.Error == nil {
			resultCh <- primary.NewTaskResult(task, result)
			return
		}
	}
	resultCh <- primary.NewTaskResult(task, secondary.SecondaryPortCallResult{Payload: nil, Error: fmt.Errorf("max retry attempt exceeded, lastError=%s", result.Error)})
}
