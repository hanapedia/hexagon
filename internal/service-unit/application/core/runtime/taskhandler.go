package runtime

import (
	"context"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
)

func TaskSetHandler(ctx context.Context, taskSet []ports.Task) []*ports.TaskError {
	done := make(chan bool, len(taskSet))
	errCh := make(chan *ports.TaskError, len(taskSet))

	for _, task := range taskSet {
		if task.Concurrent {
			go func(task ports.Task) {
				defer func() { done <- true }()
				res := task.SecondaryPort.Call(ctx)
				if res.Error != nil {
					errCh <- ports.NewTaskError(task, res.Error)
				} else {
					errCh <- nil
				}
			}(task)
		} else {
			res := task.SecondaryPort.Call(ctx)
			if res.Error != nil {
				errCh <- ports.NewTaskError(task, res.Error)
			} else {
				errCh <- nil
			}
            done <- true
		}
	}

	var errs []*ports.TaskError
	for i := 0; i < len(taskSet); i++ {
		<-done
		taskError := <-errCh
		if taskError != nil {
			errs = append(errs, taskError)
		}
	}
	return errs
}
