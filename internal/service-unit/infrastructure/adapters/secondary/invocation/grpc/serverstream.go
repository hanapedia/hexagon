package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	pb "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"google.golang.org/grpc"
)

type serverStreamAdapter struct {
	route       string
	client      *grpc.ClientConn
	payloadSize int64
	ports.SecondaryPortBase
}

func (ssa *serverStreamAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	payload, err := utils.GenerateRandomString(ssa.payloadSize)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	client := pb.NewGrpcClient(ssa.client)

	request := pb.StreamRequest{
		Route:   ssa.route,
		Message: fmt.Sprintf("Sending %v bytes to %s", ssa.payloadSize, ssa.GetDestId()),
		Payload: payload,
	}

	// server stream
	serverStream, err := client.ServerStreaming(ctx, &request)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	var lastPayload string
	for {
		resp, err := serverStream.Recv()
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
