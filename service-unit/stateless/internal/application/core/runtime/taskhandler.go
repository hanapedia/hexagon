package runtime

import (
	"context"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
)

func TaskSetHandler(ctx context.Context, taskSets []ports.TaskSet) []ports.SecondaryPortError {
	done := make(chan bool, len(taskSets))
	errCh := make(chan ports.SecondaryPortError, len(taskSets))

	for _, task := range taskSets {
		if task.Concurrent {
			go func(task ports.TaskSet) {
				defer func() { done <- true }()
				_, err := task.SecondaryPort.Call(ctx)
				errCh <- ports.SecondaryPortError{SecondaryPort: &task.SecondaryPort, Error: err}
			}(task)
		} else {
			_, err := task.SecondaryPort.Call(ctx)
			errCh <- ports.SecondaryPortError{SecondaryPort: &task.SecondaryPort, Error: err}
            done <- true
		}
	}

	var secondaryAdapterErrors []ports.SecondaryPortError
	for i := 0; i < len(taskSets); i++ {
		<-done
		invocationAdapterError := <-errCh
		if invocationAdapterError.Error != nil {
			secondaryAdapterErrors = append(secondaryAdapterErrors, invocationAdapterError)
		}
	}
	return secondaryAdapterErrors
}
