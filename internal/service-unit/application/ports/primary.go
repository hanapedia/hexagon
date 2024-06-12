package ports

import (
	"context"
	"sync"
	"time"

	model "github.com/hanapedia/hexagon/pkg/api/v1"
	l "github.com/hanapedia/hexagon/pkg/operator/logger"
)

var (
	DEFAULT_TIMEOUT = 5 * time.Second
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
	SecondaryPort SecodaryPort
	Concurrent    bool
	OnError       model.OnErrorSpec
	Timeout       string
}

// TaskResult is returned for each individual task calls
type TaskResult struct {
	Task  Task
	SecondaryPortCallResult
}

// Get parsed timeout as time.Duration
func (t Task) GetTimeout() time.Duration {
	duration, err := time.ParseDuration(t.Timeout)
	if err != nil {
		return DEFAULT_TIMEOUT
	}
	if duration == 0 {
		return DEFAULT_TIMEOUT
	}
	return duration
}


// TaskSetResult is returned for collection of task results
// It may or may not fail the request based on error handling configuration.
type TaskSetResult struct {
	ShouldFail bool
	TaskResults []*TaskResult
}

func NewTaskResult(task Task, result SecondaryPortCallResult) *TaskResult {
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

func (iah PrimaryHandler) LogTaskResult(ctx context.Context, taskResult *TaskResult) {
	if taskResult.Error != nil {
		l.Logger.WithContext(ctx).Error(
			"Call failed. ",
			"sourceId=", iah.GetId(), ", ",
			"destId=", taskResult.Task.SecondaryPort.GetDestId(), ", ",
			"err=", taskResult.SecondaryPortCallResult.Error,
		)
	} else {
		l.Logger.WithContext(ctx).Debug(
			"Call succeeded. ",
			"sourceId=", iah.GetId(), ", ",
			"destId=", taskResult.Task.SecondaryPort.GetDestId(), ", ",
		)
	}
}
