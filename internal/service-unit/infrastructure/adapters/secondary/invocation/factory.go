package invocation

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/invocation/grpc"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/invocation/rest"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func NewSecondaryAdapter(adapterConfig *model.InvocationConfig, client secondary.SecondaryAdapterClient) (secondary.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.REST:
		return rest.RestInvocationAdapterFactory(adapterConfig, client)
	case constants.GRPC:
		return grpc.GrpcInvocationAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found")
		return nil, err
	}
}

func NewClient(adapterConfig *model.InvocationConfig) secondary.SecondaryAdapterClient {
	switch adapterConfig.Variant {
	case constants.REST:
		client := rest.NewRestClient()
		return client
	case constants.GRPC:
		addr := config.GetGrpcDialAddr(adapterConfig)
		client := grpc.NewGrpcClient(addr)
		return client
	default:
		logger.Logger.Fatalf("invalid protocol")
		return nil
	}
}
