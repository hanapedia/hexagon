package secondary

import (
	"errors"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/invocation"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/producer"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/repository"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	l "github.com/hanapedia/the-bench/pkg/operator/logger"
)

func NewSecondaryAdapter(config *model.SecondaryAdapterConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	if config.InvocationConfig != nil {
		return invocation.NewSecondaryAdapter(config.InvocationConfig, client)
	}
	if config.RepositoryConfig != nil {
		return repository.NewSecondaryAdapter(config.RepositoryConfig, client)
	}
	if config.ProducerConfig != nil {
		return producer.NewSecondaryAdapter(config.ProducerConfig, client)
	}
	err := errors.New("No matching protocol found when making secondary adapter.")

	return nil, err
}

func NewSecondaryAdapterClient(config *model.SecondaryAdapterConfig) ports.SecondaryAdapter {
	if config.InvocationConfig != nil {
		return invocation.NewClient(config.InvocationConfig)
	}
	if config.RepositoryConfig != nil {
		return repository.NewClient(config.RepositoryConfig)
	}
	if config.ProducerConfig != nil {
		return producer.NewClient(config.ProducerConfig)
	}
	l.Logger.Fatalf("invalid protocol")
	return nil
}
