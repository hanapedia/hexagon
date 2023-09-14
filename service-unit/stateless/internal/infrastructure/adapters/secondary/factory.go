package secondary

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/invocation"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/repository"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/producer"
)

func NewSecondaryAdapter(config model.SecondaryAdapterConfig, clients *map[string]ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	if config.InvocationConfig != nil {
		client := invocation.GetOrCreateClient(*config.InvocationConfig, clients)
		return invocation.NewSecondaryAdapter(*config.InvocationConfig, client)
	}
	if config.RepositoryConfig != nil {
		client := repository.GetOrCreateClient(*config.RepositoryConfig, clients)
		return repository.NewSecondaryAdapter(*config.RepositoryConfig, client)
	}
	if config.ProducerConfig != nil {
		client := producer.GetOrCreateClient(*config.ProducerConfig, clients)
		return producer.NewSecondaryAdapter(*config.ProducerConfig, client)
	}
	err := errors.New("No matching protocol found when making secondary adapter.")

	return nil, err
}
