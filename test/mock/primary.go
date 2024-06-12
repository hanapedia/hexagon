package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
)

type PrimaryAdapterMock struct {
	addr string
}

// Serve mock implementation
func (pam PrimaryAdapterMock) Serve(ctx context.Context, wg *sync.WaitGroup) error {
	return nil
}

// Register mock implementation
func (pam PrimaryAdapterMock) Register(primaryHander *ports.PrimaryHandler) error {
	return nil
}

// NewPrimaryHandler returns mocked ports.PrimaryHandler with given number of tasks
func NewPrimaryHandler(numTask int) ports.PrimaryHandler {
	tasks := make([]ports.Task, numTask)
	for i := 0; i < numTask; i++ {
		tasks = append(tasks, ports.Task{
			SecondaryPort: NewSecondaryAdapter(fmt.Sprintf("task%v", i), time.Millisecond, 0),
			Concurrent:    false,
		})
	}
	return ports.PrimaryHandler{
		TaskSet: tasks,
	}
}

// NewPrimaryAdapter returns mocked ports.PrimaryPort
func NewPrimaryAdapter() ports.PrimaryPort {
	return PrimaryAdapterMock{addr: "localhost"}
}
