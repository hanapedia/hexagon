package secondary

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/invocation"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/repository"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/producer"
)

func NewSecondaryAdapter(config model.EgressAdapterConfig, clients *map[string]ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	if config.StatelessEgressAdapterConfig != nil {
		client := invocation.GetOrCreateClient(*config.StatelessEgressAdapterConfig, clients)
		return invocation.NewSecondaryAdapter(*config.StatelessEgressAdapterConfig, client)
	}
	if config.StatefulEgressAdapterConfig != nil {
		client := repository.GetOrCreateClient(*config.StatefulEgressAdapterConfig, clients)
		return repository.NewSecondaryAdapter(*config.StatefulEgressAdapterConfig, client)
	}
	if config.BrokerEgressAdapterConfig != nil {
		client := producer.GetOrCreateClient(*config.BrokerEgressAdapterConfig, clients)
		return producer.NewSecondaryAdapter(*config.BrokerEgressAdapterConfig, client)
	}
	err := errors.New("No matching protocol found when making egress adapter.")

	return nil, err
}
