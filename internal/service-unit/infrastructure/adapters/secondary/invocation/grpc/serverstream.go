package grpc

import (
	"context"
	"io"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	pb "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"google.golang.org/grpc"
)

type serverStreamAdapter struct {
	route       string
	client      *grpc.ClientConn
	payloadSize int64
	secondary.SecondaryPortBase
}

func (ssa *serverStreamAdapter) Call(ctx context.Context) secondary.SecondaryPortCallResult {
	client := pb.NewGrpcClient(ssa.client)
	payload := utils.GenerateRandomString(ssa.payloadSize)

	request := pb.StreamRequest{
		Route:   ssa.route,
		Payload: payload,
	}

	logger.Logger.Debugf("Sending request with %v bytes to %s", ssa.payloadSize, ssa.GetDestId())

	// server stream
	serverStream, err := client.ServerStreaming(ctx, &request)
	if err != nil {
		return secondary.SecondaryPortCallResult{
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
			return secondary.SecondaryPortCallResult{
				Payload: nil,
				Error:   err,
			}
		}
		lastPayload = resp.Payload
	}

	return secondary.SecondaryPortCallResult{
		Payload: &lastPayload,
		Error:   nil,
	}
}
