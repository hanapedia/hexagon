package common

import "github.com/hanapedia/the-bench/service-unit/internal/domain/core"

func TaskSetHandler(taskSets []core.TaskSet) []core.InvocationAdapterError {
	done := make(chan bool, len(taskSets))
	errCh := make(chan core.InvocationAdapterError, len(taskSets))

	for _, task := range taskSets {
		if task.Concurrent {
			go func(task core.TaskSet) {
				defer func() { done <- true }()
				_, err := task.InvocationAdapter.Call()
				errCh <- core.InvocationAdapterError{InvocationAdapter: &task.InvocationAdapter, Error: err}
			}(task)
		} else {
			_, err := task.InvocationAdapter.Call()
			errCh <- core.InvocationAdapterError{InvocationAdapter: &task.InvocationAdapter, Error: err}
            done <- true
		}
	}

	var invocationAdapterErrors []core.InvocationAdapterError
	for i := 0; i < len(taskSets); i++ {
		<-done
		invocationAdapterError := <-errCh
		if invocationAdapterError.Error != nil {
			invocationAdapterErrors = append(invocationAdapterErrors, invocationAdapterError)
		}
	}
	return invocationAdapterErrors
}
