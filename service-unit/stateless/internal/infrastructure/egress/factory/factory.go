package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
)

func NewEgressAdapter(adapterConfig model.Adapter, connections *map[string]core.EgressConnection) (core.EgressAdapter, error) {
	if adapterConfig.StatelessEgressConfig != nil {
		return statelesEgressAdapterFactory(*adapterConfig.StatelessEgressConfig)
	}
	if adapterConfig.StatefulEgressConfig != nil {
		connection := upsertStatefulEgressConnection(*adapterConfig.StatefulEgressConfig, connections)
		return statefulEgressAdapterFactory(*adapterConfig.StatefulEgressConfig, connection)
	}
	if adapterConfig.BrokerEgressConfig != nil {
		connection := upsertBrokerEgressConnection(*adapterConfig.BrokerEgressConfig, connections)
		return brokerEgressAdapterFactory(*adapterConfig.BrokerEgressConfig, connection)
	}
	err := errors.New("No matching protocol found")

	return nil, err
}
