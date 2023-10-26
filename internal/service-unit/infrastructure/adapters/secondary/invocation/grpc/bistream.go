package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	pb "github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/service-unit/payload"
	"google.golang.org/grpc"
)

type biStreamAdapter struct {
	route        string
	client       *grpc.ClientConn
	payload      constants.PayloadSizeVariant
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
		payload, err := payload.GeneratePayload(bsa.payload)
		if err != nil {
			return ports.SecondaryPortCallResult{
				Payload: nil,
				Error:   err,
			}
		}

		request := pb.StreamRequest{
			Route:   bsa.route,
			Message: fmt.Sprintf("Posting %s payload of random text to %s", bsa.payload, bsa.GetDestId()),
			Payload: payload,
		}

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
