package grpc

import (
	"errors"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
)

func GrpcInvocationAdapterFactory(adapterConfig *model.InvocationConfig, client ports.SecondaryAdapterClient) (ports.SecodaryPort, error) {
	var grpcAdapter ports.SecodaryPort
	var err error

	if grpcClient, ok := (client).(*GrpcClient); ok {
		switch adapterConfig.Action {
		case constants.SIMPLE_RPC:
			grpcAdapter = &simpleRpcAdapter{
				route:   adapterConfig.Route,
				client:  grpcClient.Connection,
				payload: adapterConfig.Payload,
			}
		case constants.SERVER_STREAM:
			grpcAdapter = &serverStreamAdapter{
				route:   adapterConfig.Route,
				client:  grpcClient.Connection,
				payload: adapterConfig.Payload,
			}
		case constants.CLIENT_STREAM:
			grpcAdapter = &clientStreamAdapter{
				route:        adapterConfig.Route,
				client:       grpcClient.Connection,
				payload:      adapterConfig.Payload,
				payloadCount: adapterConfig.PayloadCount,
			}
		case constants.BI_STREAM:
			grpcAdapter = &biStreamAdapter{
				route:        adapterConfig.Route,
				client:       grpcClient.Connection,
				payload:      adapterConfig.Payload,
				payloadCount: adapterConfig.PayloadCount,
			}
		default:
			err = errors.New("No matching protocol found when creating rest secondary adapter.")
		}
	} else {
		err = errors.New("Unmatched client instance")
	}

	// set destionation id
	grpcAdapter.SetDestId(adapterConfig.GetId())

	return grpcAdapter, err
}
