package common

import "github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"

func TaskSetHandler(taskSets []core.TaskSet) []core.EgressAdapterError {
	done := make(chan bool, len(taskSets))
	errCh := make(chan core.EgressAdapterError, len(taskSets))

	for _, task := range taskSets {
		if task.Concurrent {
			go func(task core.TaskSet) {
				defer func() { done <- true }()
				_, err := task.EgressAdapter.Call()
				errCh <- core.EgressAdapterError{EgressAdapter: &task.EgressAdapter, Error: err}
			}(task)
		} else {
			_, err := task.EgressAdapter.Call()
			errCh <- core.EgressAdapterError{EgressAdapter: &task.EgressAdapter, Error: err}
            done <- true
		}
	}

	var egressAdapterErrors []core.EgressAdapterError
	for i := 0; i < len(taskSets); i++ {
		<-done
		invocationAdapterError := <-errCh
		if invocationAdapterError.Error != nil {
			egressAdapterErrors = append(egressAdapterErrors, invocationAdapterError)
		}
	}
	return egressAdapterErrors
}
