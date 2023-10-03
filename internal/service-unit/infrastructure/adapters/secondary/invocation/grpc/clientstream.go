package grpc

import (
	"context"
	"fmt"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	pb "github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/service-unit/utils"
	"google.golang.org/grpc"
)

type clientStreamAdapter struct {
	client       *grpc.ClientConn
	payload      constants.PayloadSizeVariant
	payloadCount int
	ports.SecondaryPortBase
}

func (csa *clientStreamAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	client := pb.NewGrpcClient(csa.client)

	// Client-side streaming
	clientStream, err := client.ClientStreaming(context.Background())
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	payloadCount := csa.payloadCount
	if payloadCount <= 0 {
		payloadCount = constants.DefaultPayloadCount
	}

	for i := 0; i < payloadCount; i++ {
		payload, err := utils.GeneratePayload(csa.payload)
		if err != nil {
			return ports.SecondaryPortCallResult{
				Payload: nil,
				Error:   err,
			}
		}

		request := pb.StreamRequest{
			Message: fmt.Sprintf("Posting %s payload of random text to %s", csa.payload, csa.GetDestId()),
			Payload: payload,
		}

		clientStream.Send(&request)
	}

	response, err := clientStream.CloseAndRecv()
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
