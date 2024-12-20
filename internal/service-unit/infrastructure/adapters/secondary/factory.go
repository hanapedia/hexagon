package secondary

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/invocation"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/producer"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/repository"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/stressor"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	l "github.com/hanapedia/hexagon/pkg/operator/logger"
)

func NewSecondaryAdapter(config *model.SecondaryAdapterConfig, client secondary.SecondaryAdapterClient) (secondary.SecodaryPort, error) {
	if config.InvocationConfig != nil {
		return invocation.NewSecondaryAdapter(config.InvocationConfig, client)
	}
	if config.RepositoryConfig != nil {
		return repository.NewSecondaryAdapter(config.RepositoryConfig, client)
	}
	if config.ProducerConfig != nil {
		return producer.NewSecondaryAdapter(config.ProducerConfig, client)
	}
	if config.StressorConfig != nil {
		return stressor.NewSecondaryAdapter(config.StressorConfig, client)
	}
	err := errors.New("No matching protocol found when making secondary adapter.")

	return nil, err
}

func NewSecondaryAdapterClient(config *model.SecondaryAdapterConfig) secondary.SecondaryAdapterClient {
	if config.InvocationConfig != nil {
		return invocation.NewClient(config.InvocationConfig)
	}
	if config.RepositoryConfig != nil {
		return repository.NewClient(config.RepositoryConfig)
	}
	if config.ProducerConfig != nil {
		return producer.NewClient(config.ProducerConfig)
	}
	if config.StressorConfig != nil {
		return stressor.NewClient(config.StressorConfig)
	}
	l.Logger.Fatalf("invalid protocol")
	return nil
}
