package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/primary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
)

type PrimaryAdapterMock struct {
	addr string
}

// Serve mock implementation
func (pam PrimaryAdapterMock) Serve(ctx context.Context, wg *sync.WaitGroup) error {
	return nil
}

// Register mock implementation
func (pam PrimaryAdapterMock) Register(primaryHander *domain.PrimaryAdapterHandler) error {
	return nil
}

// NewPrimaryHandler returns mocked ports.PrimaryHandler with given number of tasks
func NewPrimaryHandler(numTask int) domain.PrimaryAdapterHandler {
	tasks := make([]domain.Task, numTask)
	for i := 0; i < numTask; i++ {
		tasks = append(tasks, domain.Task{
			SecondaryPort: NewSecondaryAdapter(fmt.Sprintf("task%v", i), time.Millisecond, 0),
			Concurrent:    false,
		})
	}
	return domain.PrimaryAdapterHandler{
		TaskSet: tasks,
	}
}

// NewPrimaryAdapter returns mocked ports.PrimaryPort
func NewPrimaryAdapter() primary.PrimaryPort {
	return PrimaryAdapterMock{addr: "localhost"}
}
