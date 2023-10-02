package grpc

import (
	"context"
	"io"

	pb "github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
)

type GrpcServerAdapter struct {
	addr    string
	payload constants.PayloadSizeVariant
	pb.UnimplementedGrpcServer
}

// Regular RPC
func (s *GrpcServerAdapter) RegularCall(ctx context.Context, req *pb.StreamRequest) (*pb.StreamResponse, error) {
	return &pb.StreamResponse{Message: "Received", Payload: req.Payload}, nil
}

// Client-side streaming
func (s *GrpcServerAdapter) ClientStreaming(stream pb.Grpc_ClientStreamingServer) error {
	var lastPayload string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Message: "Received", Payload: lastPayload})
		}
		if err != nil {
			return err
		}
		lastPayload = req.Payload
	}
}

// Server-side streaming
func (s *GrpcServerAdapter) ServerStreaming(req *pb.StreamRequest, stream pb.Grpc_ServerStreamingServer) error {
	for i := 0; i < 3; i++ {
		if err := stream.Send(&pb.StreamResponse{Message: "Streaming", Payload: req.Payload}); err != nil {
			return err
		}
	}
	return nil
}

// Bidirectional streaming
func (s *GrpcServerAdapter) BidirectionalStreaming(stream pb.Grpc_BidirectionalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.StreamResponse{Message: "Streaming", Payload: req.Payload}); err != nil {
			return err
		}
	}
}
