package invocation

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/invocation/rest"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

func NewSecondaryAdapter(adapterConfig *model.InvocationConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.REST:
		return rest.RestInvocationAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found")
		return nil, err
	}
}

func NewClient(adapterConfig *model.InvocationConfig) ports.SecondaryAdapter {
	switch adapterConfig.Variant {
	case constants.REST:
		client := rest.NewRestClient()
		return client
	default:
		logger.Logger.Fatalf("invalid protocol")
		return nil
	}
}
