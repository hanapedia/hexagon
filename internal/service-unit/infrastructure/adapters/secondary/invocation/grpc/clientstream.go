package grpc

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	pb "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"google.golang.org/grpc"
)

type clientStreamAdapter struct {
	route        string
	client       *grpc.ClientConn
	payloadSize  int64
	payloadCount int
	secondary.SecondaryPortBase
}

func (csa *clientStreamAdapter) Call(ctx context.Context) secondary.SecondaryPortCallResult {
	client := pb.NewGrpcClient(csa.client)

	// Client-side streaming
	clientStream, err := client.ClientStreaming(ctx)
	if err != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	payloadCount := csa.payloadCount
	if payloadCount <= 0 {
		payloadCount = constants.DefaultPayloadCount
	}

	for i := 0; i < payloadCount; i++ {
		payload := utils.GenerateRandomString(csa.payloadSize)
		request := pb.StreamRequest{
			Route:   csa.route,
			Payload: payload,
		}

		logger.Logger.Debugf("Sending request with %v bytes to %s", csa.payloadSize, csa.GetDestId())

		clientStream.Send(&request)
	}

	response, err := clientStream.CloseAndRecv()
	if err != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return secondary.SecondaryPortCallResult{
		Payload: &response.Payload,
		Error:   nil,
	}
}
