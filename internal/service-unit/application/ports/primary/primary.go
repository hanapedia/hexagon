package primary

import (
	"context"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	l "github.com/hanapedia/hexagon/pkg/operator/logger"
)

var (
	DEFAULT_TASK_TIMEOUT = 5 * time.Second
	DEFAULT_CALL_TIMEOUT = 1 * time.Second
)

// PrimaryPort provides common interface for all the primary adapters.
// Example resources include:
// - REST API server
// - gRPC server
// - Kafka consumer
//
// It is intended to represent the individual interfaces on each exteranl service,
type PrimaryPort interface {
	// Serve starts primary port adapter with cancellable context and WaitGroup for graceful shutdown
	Serve(context.Context, *sync.WaitGroup) error
	// Register registers primary port handler to primary adapter instance
	Register(*PrimaryHandler) error
}

type PrimaryAdapterHandler interface {
	GetId() string
	GetConfig() PrimaryAdapterConfig
}

type PrimaryAdapterConfig interface {
	GetId(string) string
	GetGroupByKey() string
}

type PrimaryPortError struct {
	PrimaryPort PrimaryPort
	Error       error
}

// either StatelessAdapterConfig or BrokerAdapterConfig must be defined
type PrimaryHandler struct {
	ServiceName    string
	ServerConfig   *model.ServerConfig
	ConsumerConfig *model.ConsumerConfig
	TaskSet        []Task
}

type Task struct {
	SecondaryPort secondary.SecodaryPort
	Concurrent    bool
	OnError       model.OnErrorSpec
	TaskTimeout   string
	CallTimeout   string
}

// TaskResult is returned for each individual task calls
type TaskResult struct {
	Task Task
	secondary.SecondaryPortCallResult
}

// Get parsed taskTimeout as time.Duration
func (t Task) GetTaskTimeout() time.Duration {
	duration, err := time.ParseDuration(t.TaskTimeout)
	if err != nil {
		return DEFAULT_TASK_TIMEOUT
	}
	if duration == 0 {
		return DEFAULT_TASK_TIMEOUT
	}
	return duration
}

// Get parsed getCallTimeout as time.Duration
func (t Task) GetCallTimeout() time.Duration {
	duration, err := time.ParseDuration(t.CallTimeout)
	if err != nil {
		return DEFAULT_CALL_TIMEOUT
	}
	if duration == 0 {
		return DEFAULT_CALL_TIMEOUT
	}
	return duration
}

// TaskSetResult is returned for collection of task results
// It may or may not fail the request based on error handling configuration.
type TaskSetResult struct {
	ShouldFail  bool
	TaskResults []*TaskResult
}

func NewTaskResult(task Task, result secondary.SecondaryPortCallResult) *TaskResult {
	return &TaskResult{Task: task, SecondaryPortCallResult: result}
}

func (iah PrimaryHandler) GetId() string {
	var id string
	if iah.ServerConfig != nil {
		id = iah.ServerConfig.GetId(iah.ServiceName)
	}
	if iah.ConsumerConfig != nil {
		id = iah.ConsumerConfig.GetId(iah.ServiceName)
	}
	return id
}

func LogTaskResult(ctx context.Context, adapterId string, taskResult *TaskResult) {
	if taskResult.Error != nil {
		l.Logger.WithContext(ctx).Error(
			"Call failed. ",
			"sourceId=", adapterId, ", ",
			"destId=", taskResult.Task.SecondaryPort.GetDestId(), ", ",
			"err=", taskResult.SecondaryPortCallResult.Error,
		)
	} else {
		l.Logger.WithContext(ctx).Debug(
			"Call succeeded. ",
			"sourceId=", adapterId, ", ",
			"destId=", taskResult.Task.SecondaryPort.GetDestId(), ", ",
		)
	}
}
