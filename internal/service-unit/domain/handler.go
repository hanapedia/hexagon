package domain

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

type TaskHandler func(context.Context, chan<- *secondary.SecondaryPortCallResult)

// either StatelessAdapterConfig or BrokerAdapterConfig must be defined
type PrimaryAdapterHandler struct {
	ServiceName    string
	ServerConfig   *model.ServerConfig
	ConsumerConfig *model.ConsumerConfig
	TaskSet        []TaskHandler
}


func (iah PrimaryAdapterHandler) GetId() string {
	var id string
	if iah.ServerConfig != nil {
		id = iah.ServerConfig.GetId(iah.ServiceName)
	}
	if iah.ConsumerConfig != nil {
		id = iah.ConsumerConfig.GetId(iah.ServiceName)
	}
	return id
}

// TaskSetResult is returned for collection of task results
// It may or may not fail the request based on error handling configuration.
type TaskSetResult struct {
	ShouldFail  bool
	TaskResults []*secondary.SecondaryPortCallResult
}
