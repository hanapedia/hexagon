package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	pb "github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/service-unit/utils"
	"google.golang.org/grpc"
)

type serverStreamAdapter struct {
	route   string
	client  *grpc.ClientConn
	payload constants.PayloadSizeVariant
	ports.SecondaryPortBase
}

func (ssa *serverStreamAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	payload, err := utils.GeneratePayload(ssa.payload)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	client := pb.NewGrpcClient(ssa.client)

	request := pb.StreamRequest{
		Route:   ssa.route,
		Message: fmt.Sprintf("Posting %s payload of random text to %s", ssa.payload, ssa.GetDestId()),
		Payload: payload,
	}

	// server stream
	serverStream, err := client.ServerStreaming(context.Background(), &request)
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
			lastPayload = resp.Payload
			break
		}
		if err != nil {
			return ports.SecondaryPortCallResult{
				Payload: nil,
				Error:   err,
			}
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: &lastPayload,
		Error:   nil,
	}
}
