package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/hanapedia/the-bench/internal/service-unit/application/core/runtime"
	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	pb "github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/service-unit/utils"
	"google.golang.org/grpc"
)

// must implement ports.PrimaryPort
type GrpcServerAdapter struct {
	addr    string
	server  *grpc.Server
	configs GrpcVariantConfigs
	pb.UnimplementedGrpcServer
}

// GrpcVariantConfigs holds the config for each grpc variant
type GrpcVariantConfigs struct {
	simpleRpc     map[string]*ports.PrimaryHandler
	clientStream  map[string]*ports.PrimaryHandler
	serverStream  map[string]*ports.PrimaryHandler
	bidirectional map[string]*ports.PrimaryHandler
}

func NewGrpcServerAdapter() *GrpcServerAdapter {
	server := grpc.NewServer()

	adapter := GrpcServerAdapter{
		addr:   config.GetRestServerAddr(),
		server: server,
	}
	return &adapter
}

// Serve starts the grpc server
func (gsa *GrpcServerAdapter) Serve() error {
	listen, err := net.Listen("tcp", gsa.addr)
	if err != nil {
		return err
	}

	pb.RegisterGrpcServer(gsa.server, &GrpcServerAdapter{})

	return gsa.server.Serve(listen)
}

// Register registers tasks to the server
func (gsa *GrpcServerAdapter) Register(handler *ports.PrimaryHandler) error {
	if handler.ServerConfig == nil {
		return errors.New(fmt.Sprintf("Invalid configuartion for handler %s.", handler.GetId()))
	}

	var simpleRpc map[string]*ports.PrimaryHandler
	var clientStream map[string]*ports.PrimaryHandler
	var serverStream map[string]*ports.PrimaryHandler
	var bidirectional map[string]*ports.PrimaryHandler

	switch handler.ServerConfig.Action {
	case constants.RPC:
		simpleRpc[handler.ServerConfig.Route] = handler
	case constants.CLIENT_STREAM:
		clientStream[handler.ServerConfig.Route] = handler
	case constants.SERVER_STREAM:
		serverStream[handler.ServerConfig.Route] = handler
	case constants.BIDIRECTIONAL:
		bidirectional[handler.ServerConfig.Route] = handler
	}

	gsa.configs = GrpcVariantConfigs{
		simpleRpc:     simpleRpc,
		clientStream:  clientStream,
		serverStream:  serverStream,
		bidirectional: bidirectional,
	}

	return nil
}

// Regular RPC
func (s *GrpcServerAdapter) SimpleRPC(ctx context.Context, req *pb.StreamRequest) (*pb.StreamResponse, error) {
	handler, ok := s.configs.simpleRpc[req.Route]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Route not found %s.", req.Route))
	}

	errs := runtime.TaskSetHandler(ctx, handler.TaskSet)
	if errs != nil {
		for _, err := range errs {
			handler.LogTaskError(ctx, err)
		}
		return nil, errors.New(fmt.Sprintf("Simple RPC %s failed when handling tasks.", handler.GetId()))
	}

	// write response
	payload, err := utils.GeneratePayload(handler.ServerConfig.Payload)
	if err != nil {
		return nil, err
	}

	rpcResponse := pb.StreamResponse{
		Message: fmt.Sprintf("Successfully ran %s, sending %s payload.", handler.GetId(), handler.ServerConfig.Payload),
		Payload: payload,
	}

	return &rpcResponse, nil
}

// Client-side streaming
func (s *GrpcServerAdapter) ClientStreaming(stream pb.Grpc_ClientStreamingServer) error {
	// process the first message in the stream and start tasks
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	handler, ok := s.configs.clientStream[req.Route]
	if !ok {
		return errors.New(fmt.Sprintf("Route not found %s.", req.Route))
	}

	errs := runtime.TaskSetHandler(stream.Context(), handler.TaskSet)
	if errs != nil {
		for _, err := range errs {
			handler.LogTaskError(stream.Context(), err)
		}
		return errors.New(fmt.Sprintf("Client streaming %s failed when handling tasks.", handler.GetId()))
	}

	// process rest of the stream
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			// write response
			payload, err := utils.GeneratePayload(handler.ServerConfig.Payload)
			if err != nil {
				return err
			}

			rpcResponse := pb.StreamResponse{
				Message: fmt.Sprintf("Successfully ran %s, sending %s payload.", handler.GetId(), handler.ServerConfig.Payload),
				Payload: payload,
			}
			return stream.SendAndClose(&rpcResponse)
		}
		if err != nil {
			return err
		}
	}
}

// Server-side streaming
func (s *GrpcServerAdapter) ServerStreaming(req *pb.StreamRequest, stream pb.Grpc_ServerStreamingServer) error {
	handler, ok := s.configs.serverStream[req.Route]
	if !ok {
		return errors.New(fmt.Sprintf("Route not found %s.", req.Route))
	}

	errs := runtime.TaskSetHandler(stream.Context(), handler.TaskSet)
	if errs != nil {
		for _, err := range errs {
			handler.LogTaskError(stream.Context(), err)
		}
		return errors.New(fmt.Sprintf("Server streaming %s failed when handling tasks.", handler.GetId()))
	}

	payloadCount := handler.ServerConfig.PayloadCount
	if payloadCount <= 0 {
		payloadCount = constants.DefaultPayloadCount
	}

	for i := 0; i < payloadCount; i++ {
		payload, err := utils.GeneratePayload(handler.ServerConfig.Payload)
		if err != nil {
			return err
		}

		rpcResponse := pb.StreamResponse{
			Message: fmt.Sprintf("Successfully ran %s, sending %s payload.", handler.GetId(), handler.ServerConfig.Payload),
			Payload: payload,
		}
		if err := stream.Send(&rpcResponse); err != nil {
			return err
		}
	}
	return nil
}

// Bidirectional streaming
func (s *GrpcServerAdapter) BidirectionalStreaming(stream pb.Grpc_BidirectionalStreamingServer) error {
	// process the first message in the stream and start tasks
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	handler, ok := s.configs.serverStream[req.Route]
	if !ok {
		return errors.New(fmt.Sprintf("Route not found %s.", req.Route))
	}

	errs := runtime.TaskSetHandler(stream.Context(), handler.TaskSet)
	if errs != nil {
		for _, err := range errs {
			handler.LogTaskError(stream.Context(), err)
		}
		return errors.New(fmt.Sprintf("Bidirectional streaming %s failed when handling tasks.", handler.GetId()))
	}

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		payload, err := utils.GeneratePayload(handler.ServerConfig.Payload)
		if err != nil {
			return err
		}

		rpcResponse := pb.StreamResponse{
			Message: fmt.Sprintf("Successfully ran %s, sending %s payload.", handler.GetId(), handler.ServerConfig.Payload),
			Payload: payload,
		}
		if err := stream.Send(&rpcResponse); err != nil {
			return err
		}
	}
}
