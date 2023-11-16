package grpc

import (
	"context"
	"io"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	pb "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"google.golang.org/grpc"
)

type biStreamAdapter struct {
	route        string
	client       *grpc.ClientConn
	payloadSize  int64
	payloadCount int
	ports.SecondaryPortBase
}

func (bsa *biStreamAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	client := pb.NewGrpcClient(bsa.client)

	// Client-side streaming
	biStream, err := client.BidirectionalStreaming(ctx)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	payloadCount := bsa.payloadCount
	if payloadCount <= 0 {
		payloadCount = constants.DefaultPayloadCount
	}

	for i := 0; i < payloadCount; i++ {
		payload := utils.GenerateRandomString(bsa.payloadSize)
		request := pb.StreamRequest{
			Route:   bsa.route,
			Payload: payload,
		}

		logger.Logger.Debugf("Sending request with %v bytes to %s", bsa.payloadSize, bsa.GetDestId())

		biStream.Send(&request)
	}
	biStream.CloseSend()

	var lastPayload string
	for {
		resp, err := biStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return ports.SecondaryPortCallResult{
				Payload: nil,
				Error:   err,
			}
		}
		lastPayload = resp.Payload
	}

	return ports.SecondaryPortCallResult{
		Payload: &lastPayload,
		Error:   nil,
	}
}
