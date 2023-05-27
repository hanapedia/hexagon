package factory

import (
	"errors"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
)

func NewEgressAdapter(egressAdapterConfig model.EgressAdapterConfig, client *map[string]core.EgressClient) (core.EgressAdapter, error) {
	if egressAdapterConfig.StatelessEgressAdapterConfig != nil {
		client := getOrCreateStatelessEgressClient(*egressAdapterConfig.StatelessEgressAdapterConfig, client)
		return statelesEgressAdapterFactory(*egressAdapterConfig.StatelessEgressAdapterConfig, client)
	}
	if egressAdapterConfig.StatefulEgressAdapterConfig != nil {
		client := getOrCreateStatefulEgressClient(*egressAdapterConfig.StatefulEgressAdapterConfig, client)
		return statefulEgressAdapterFactory(*egressAdapterConfig.StatefulEgressAdapterConfig, client)
	}
	if egressAdapterConfig.BrokerEgressAdapterConfig != nil {
		client := getOrCreateBrokerEgressClient(*egressAdapterConfig.BrokerEgressAdapterConfig, client)
		return brokerEgressAdapterFactory(*egressAdapterConfig.BrokerEgressAdapterConfig, client)
	}
	err := errors.New("No matching protocol found when making egress adapter.")

	return nil, err
}
