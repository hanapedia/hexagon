package ports

import (
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
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
	PrimaryPort *PrimaryPort
	Error       error
}

// either StatelessAdapterConfig or BrokerAdapterConfig must be defined
type PrimaryHandler struct {
	ServerConfig   *model.ServerConfig
	ConsumerConfig *model.ConsumerConfig
	TaskSets       []TaskSet
}

type TaskSet struct {
	SecondaryPort SecodaryPort
	Concurrent    bool
}

func (iah PrimaryHandler) GetId(serviceName string) string {
	var id string
	if iah.ServerConfig != nil {
		id = iah.ServerConfig.GetId(serviceName)
	}
	if iah.ConsumerConfig != nil {
		id = iah.ConsumerConfig.GetId(serviceName)
	}
	return id
}
