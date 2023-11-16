package grpc

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	pb "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"google.golang.org/grpc"
)

type simpleRpcAdapter struct {
	route       string
	client      *grpc.ClientConn
	payloadSize int64
	ports.SecondaryPortBase
}

func (sra *simpleRpcAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	payload, err := utils.GenerateRandomString(sra.payloadSize)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	client := pb.NewGrpcClient(sra.client)

	request := pb.StreamRequest{
		Route:   sra.route,
		Payload: payload,
	}

	logger.Logger.Debugf("Sending request with %v bytes to %s", sra.payloadSize, sra.GetDestId())

	// Regular RPC
	response, err := client.SimpleRPC(ctx, &request)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: &response.Payload,
		Error:   nil,
	}
}
