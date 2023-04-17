package factory

import (
	"errors"
	"strings"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/invocation_adapter/kafka"
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
)

type InvocationAdapterDetails struct {
	serviceName string
	protocol    string
	action      string
	handlerName string
	connection  core.Connection
}

func newInvocationAdapterDetails(id string) (InvocationAdapterDetails, error) {
	idSubstring := strings.Split(id, ".")
	var err error
	if len(idSubstring) != 4 {
		err = errors.New("Invalid adapter id")
	}
	return InvocationAdapterDetails{
		serviceName: idSubstring[constants.ServiceNameIndex],
		protocol:    idSubstring[constants.ProtocolIndex],
		action:      idSubstring[constants.ActionIndex],
		handlerName: idSubstring[constants.AdapterNameIndex],
	}, err
}

func NewInvocationAdapter(id string, connections *map[string]core.Connection) (core.InvocationAdapter, error) {
	invocationAdapterDetails, err := newInvocationAdapterDetails(id)
	var invocationAdapter core.InvocationAdapter
	switch invocationAdapterDetails.protocol {
	case "rest":
		invocationAdapter, err = invocationAdapterDetails.restInvocationAdapterFactory()
	case "kafka":
        invocationAdapterDetails.UpsertConnection(id, connections)
		invocationAdapter, err = invocationAdapterDetails.kafkaInvocationAdapterFactory()
	default:
		err = errors.New("No matching protocol found")
	}

	return invocationAdapter, err
}

func (invocationAdapterDetails *InvocationAdapterDetails) UpsertConnection(id string, connections *map[string]core.Connection) {
	connection, ok := (*connections)[id]
	if ok {
		invocationAdapterDetails.connection = connection
		return
	}
	switch invocationAdapterDetails.protocol {
	case "kafka":
		invocationAdapterDetails.connection = kafka.NewKafkaConnection(constants.KafkaBrokerAddr, invocationAdapterDetails.action)
	}
}
