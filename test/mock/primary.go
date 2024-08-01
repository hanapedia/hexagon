package mock

import (
	"context"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/primary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/resiliency"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
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
	tasks := make([]domain.TaskHandler, numTask)
	for i := 0; i < numTask; i++ {
		tasks = append(tasks,
			resiliency.NewTaskHandler(
				domain.TelemetryContext{PrimaryLabels: domain.PrimaryLabels{ServiceName: "RegularCallHandler"}},
				model.TaskSpec{},
				NewSecondaryAdapter("RegularSecondaryAdapter1", time.Second, 0),
			))
	}
	return domain.PrimaryAdapterHandler{
		TaskSet: tasks,
	}
}

// NewPrimaryAdapter returns mocked ports.PrimaryPort
func NewPrimaryAdapter() primary.PrimaryPort {
	return PrimaryAdapterMock{addr: "localhost"}
}
