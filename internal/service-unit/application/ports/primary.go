package ports

import (
	"context"

	model "github.com/hanapedia/the-bench/pkg/api/v1"
	l "github.com/hanapedia/the-bench/pkg/operator/logger"
)

// PrimaryPort provides common interface for all the primary adapters.
// Example resources include:
// - REST API server
// - gRPC server
// - Kafka consumer
//
// It is intended to represent the individual interfaces on each exteranl service,
type PrimaryPort interface {
	Serve() error
	Register(string, *PrimaryHandler) error
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
}

type TaskError struct {
	task  Task
	error error
}

func NewTaskError(task Task, err error) *TaskError {
	return &TaskError{task: task, error: err}
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

func (iah PrimaryHandler) LogTaskError(ctx context.Context, taskError *TaskError) {
	l.Logger.WithContext(ctx).Error(
		"Call failed. ",
		"sourceId=", iah.GetId(), ", ",
		"destId=", taskError.task.SecondaryPort.GetDestId(), ", ",
		"err=", taskError.error,
	)
}
