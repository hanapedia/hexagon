package factory

import (
	"errors"
	"strings"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/producer_adapter/kafka"
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
)

type EgressAdapterDetails struct {
	serviceName string
	protocol    constants.AdapterProtocol
	action      string
	handlerName string
	connection  core.Connection
}

func newEgressAdapterDetails(id string) (EgressAdapterDetails, error) {
	idSubstring := strings.Split(id, ".")
	var err error
	if len(idSubstring) != 4 {
		err = errors.New("Invalid adapter id")
	}
	return EgressAdapterDetails{
		serviceName: idSubstring[constants.ServiceNameIndex],
		protocol:    constants.AdapterProtocol(idSubstring[constants.ProtocolIndex]),
		action:      idSubstring[constants.ActionIndex],
		handlerName: idSubstring[constants.AdapterNameIndex],
	}, err
}

func NewEgressAdapter(id string, connections *map[string]core.Connection) (core.EgressAdapter, error) {
	egressAdapterDetails, err := newEgressAdapterDetails(id)
	var egressAdapter core.EgressAdapter
	switch egressAdapterDetails.protocol {
	case constants.REST:
		egressAdapter, err = egressAdapterDetails.restEgressAdapterFactory()
	case constants.KAFKA:
        egressAdapterDetails.UpsertConnection(id, connections)
		egressAdapter, err = egressAdapterDetails.kafkaEgressAdapterFactory()
	default:
		err = errors.New("No matching protocol found")
	}

	return egressAdapter, err
}

func (egressAdapterDetails *EgressAdapterDetails) UpsertConnection(id string, connections *map[string]core.Connection) {
	connection, ok := (*connections)[id]
	if ok {
		egressAdapterDetails.connection = connection
		return
	}
	switch egressAdapterDetails.protocol {
	case "kafka":
		egressAdapterDetails.connection = kafka.NewKafkaConnection(constants.KafkaBrokerAddr, egressAdapterDetails.action)
	}
}
