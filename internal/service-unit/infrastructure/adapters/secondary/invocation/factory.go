package invocation

import (
	"errors"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/invocation/rest"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
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