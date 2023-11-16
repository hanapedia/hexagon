package grpc

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func GrpcInvocationAdapterFactory(adapterConfig *model.InvocationConfig, client ports.SecondaryAdapterClient) (ports.SecodaryPort, error) {
	var grpcAdapter ports.SecodaryPort
	var err error

	if grpcClient, ok := (client).(*GrpcClient); ok {
		payloadSize := model.GetPayloadSize(adapterConfig.Payload)
		switch adapterConfig.Action {
		case constants.SIMPLE_RPC:
			grpcAdapter = &simpleRpcAdapter{
				route:       adapterConfig.Route,
				client:      grpcClient.Connection,
				payloadSize: payloadSize,
			}
		case constants.SERVER_STREAM:
			grpcAdapter = &serverStreamAdapter{
				route:       adapterConfig.Route,
				client:      grpcClient.Connection,
				payloadSize: payloadSize,
			}
		case constants.CLIENT_STREAM:
			grpcAdapter = &clientStreamAdapter{
				route:        adapterConfig.Route,
				client:       grpcClient.Connection,
				payloadSize:  payloadSize,
				payloadCount: adapterConfig.Payload.Count,
			}
		case constants.BI_STREAM:
			grpcAdapter = &biStreamAdapter{
				route:        adapterConfig.Route,
				client:       grpcClient.Connection,
				payloadSize:  payloadSize,
				payloadCount: adapterConfig.Payload.Count,
			}
		default:
			err = errors.New("No matching protocol found when creating rest secondary adapter.")
		}
	} else {
		err = errors.New("Unmatched client instance")
	}

	// set destionation id
	grpcAdapter.SetDestId(adapterConfig.GetId())

	logger.Logger.Debugf("Initialized gRPC invocation adapter: %s", adapterConfig.GetId())
	return grpcAdapter, err
}
