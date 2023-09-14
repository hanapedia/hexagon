package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

func NewEgressAdapter(egressAdapterConfig model.EgressAdapterConfig, clients *map[string]ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	if egressAdapterConfig.StatelessEgressAdapterConfig != nil {
		client := getOrCreateStatelessEgressClient(*egressAdapterConfig.StatelessEgressAdapterConfig, clients)
		return statelesEgressAdapterFactory(*egressAdapterConfig.StatelessEgressAdapterConfig, client)
	}
	if egressAdapterConfig.StatefulEgressAdapterConfig != nil {
		client := getOrCreateStatefulEgressClient(*egressAdapterConfig.StatefulEgressAdapterConfig, clients)
		return statefulEgressAdapterFactory(*egressAdapterConfig.StatefulEgressAdapterConfig, client)
	}
	if egressAdapterConfig.BrokerEgressAdapterConfig != nil {
		client := getOrCreateBrokerEgressClient(*egressAdapterConfig.BrokerEgressAdapterConfig, clients)
		return brokerEgressAdapterFactory(*egressAdapterConfig.BrokerEgressAdapterConfig, client)
	}
	err := errors.New("No matching protocol found when making egress adapter.")

	return nil, err
}
