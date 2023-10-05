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

type simpleRpcAdapter struct {
	route   string
	client  *grpc.ClientConn
	payload constants.PayloadSizeVariant
	ports.SecondaryPortBase
}

func (sra *simpleRpcAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	payload, err := utils.GeneratePayload(sra.payload)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	client := pb.NewGrpcClient(sra.client)

	request := pb.StreamRequest{
		Route:   sra.route,
		Message: fmt.Sprintf("Posting %s payload of random text to %s", sra.payload, sra.GetDestId()),
		Payload: payload,
	}

	// Regular RPC
	response, err := client.SimpleRPC(context.Background(), &request)
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
