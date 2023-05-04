package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
)

func NewEgressAdapter(egressAdapterConfig model.EgressAdapterConfig, connections *map[string]core.EgressConnection) (core.EgressAdapter, error) {
	if egressAdapterConfig.StatelessEgressAdapterConfig != nil {
		return statelesEgressAdapterFactory(*egressAdapterConfig.StatelessEgressAdapterConfig)
	}
	if egressAdapterConfig.StatefulEgressAdapterConfig != nil {
		connection := upsertStatefulEgressConnection(*egressAdapterConfig.StatefulEgressAdapterConfig, connections)
		return statefulEgressAdapterFactory(*egressAdapterConfig.StatefulEgressAdapterConfig, connection)
	}
	if egressAdapterConfig.BrokerEgressAdapterConfig != nil {
		connection := upsertBrokerEgressConnection(*egressAdapterConfig.BrokerEgressAdapterConfig, connections)
		return brokerEgressAdapterFactory(*egressAdapterConfig.BrokerEgressAdapterConfig, connection)
	}
	err := errors.New("No matching protocol found when making egress adapter.")

	return nil, err
}
