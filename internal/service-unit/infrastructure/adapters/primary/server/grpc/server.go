package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/core/runtime"
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	pb "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/generated/grpc"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	util "github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// must implement ports.PrimaryPort
type GrpcServerAdapter struct {
	addr        string
	server      *grpc.Server
	configs     GrpcVariantConfigs
	pb.UnimplementedGrpcServer
}

// GrpcVariantConfigs holds the config for each grpc variant
type GrpcVariantConfigs struct {
	simpleRpc    map[string]*ports.PrimaryHandler
	clientStream map[string]*ports.PrimaryHandler
	serverStream map[string]*ports.PrimaryHandler
	biStream     map[string]*ports.PrimaryHandler
}

func NewGrpcServerAdapter() *GrpcServerAdapter {
	var opts []grpc.ServerOption

	// enable tracing
	if config.GetEnvs().TRACING {
		opts = append(opts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}

	server := grpc.NewServer(opts...)

	adapter := GrpcServerAdapter{
		addr:   config.GetGrpcServerAddr(),
		server: server,
		configs: GrpcVariantConfigs{
			simpleRpc:    make(map[string]*ports.PrimaryHandler),
			clientStream: make(map[string]*ports.PrimaryHandler),
			serverStream: make(map[string]*ports.PrimaryHandler),
			biStream:     make(map[string]*ports.PrimaryHandler),
		},
	}
	return &adapter
}

// Serve starts the grpc server
func (gsa *GrpcServerAdapter) Serve(ctx context.Context, wg *sync.WaitGroup) error {
	logger.Logger.Infof("Serving grpc server at %s", gsa.addr)
	go func() {
		<- ctx.Done()
		logger.Logger.Infof("Context cancelled. GRPC Server shutting down.")
		gsa.server.GracefulStop()
		wg.Done()
	}()

	listen, err := net.Listen("tcp", gsa.addr)
	if err != nil {
		return err
	}

	pb.RegisterGrpcServer(gsa.server, gsa)

	return gsa.server.Serve(listen)
}

// Register registers tasks to the server
func (gsa *GrpcServerAdapter) Register(handler *ports.PrimaryHandler) error {
	if handler.ServerConfig == nil {
		return errors.New(fmt.Sprintf("Invalid configuartion for handler %s.", handler.GetId()))
	}

	switch handler.ServerConfig.Action {
	case constants.SIMPLE_RPC:
		logger.Logger.Debugf("Registered simple rpc at %s", handler.ServerConfig.Route)
		gsa.configs.simpleRpc[handler.ServerConfig.Route] = handler
	case constants.CLIENT_STREAM:
		logger.Logger.Debugf("Registered client stream at %s", handler.ServerConfig.Route)
		gsa.configs.clientStream[handler.ServerConfig.Route] = handler
	case constants.SERVER_STREAM:
		logger.Logger.Debugf("Registered server stream at %s", handler.ServerConfig.Route)
		gsa.configs.serverStream[handler.ServerConfig.Route] = handler
	case constants.BI_STREAM:
		logger.Logger.Debugf("Registered bi stream at %s", handler.ServerConfig.Route)
		gsa.configs.biStream[handler.ServerConfig.Route] = handler
	}
	return nil
}

// Regular RPC
func (gsa *GrpcServerAdapter) SimpleRPC(ctx context.Context, req *pb.StreamRequest) (*pb.StreamResponse, error) {
	// record time for logging
	startTime := time.Now()

	handler, ok := gsa.configs.simpleRpc[req.Route]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Route not found %s.", req.Route))
	}

	// defer log
	defer gsa.log(ctx, handler, startTime)

	errs := runtime.TaskSetHandler(ctx, handler.TaskSet)
	if errs != nil {
		for _, err := range errs {
			handler.LogTaskError(ctx, err)
		}
		return nil, errors.New(fmt.Sprintf("Simple RPC %s failed when handling tasks.", handler.GetId()))
	}

	// write response
	payloadSize := model.GetPayloadSize(handler.ServerConfig.Payload)
	payload := util.GenerateRandomString(payloadSize)

	rpcResponse := pb.StreamResponse{
		Payload: payload,
	}

	logger.Logger.Debugf("Ran %s, responding with %v bytes.", handler.GetId(), payloadSize)

	return &rpcResponse, nil
}

// Client-side streaming
func (gsa *GrpcServerAdapter) ClientStreaming(stream pb.Grpc_ClientStreamingServer) error {
	// record time for logging
	startTime := time.Now()

	// process the first message in the stream and start tasks
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	handler, ok := gsa.configs.clientStream[req.Route]
	if !ok {
		return errors.New(fmt.Sprintf("Route not found %s.", req.Route))
	}

	// defer log
	defer gsa.log(stream.Context(), handler, startTime)

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
			payloadSize := model.GetPayloadSize(handler.ServerConfig.Payload)
			payload := util.GenerateRandomString(payloadSize)

			rpcResponse := pb.StreamResponse{
				Payload: payload,
			}
			logger.Logger.Debugf("Ran %s, responding with %v bytes.", handler.GetId(), payloadSize)
			return stream.SendAndClose(&rpcResponse)
		}
		if err != nil {
			return err
		}
	}
}

// Server-side streaming
func (gsa *GrpcServerAdapter) ServerStreaming(req *pb.StreamRequest, stream pb.Grpc_ServerStreamingServer) error {
	// record time for logging
	startTime := time.Now()

	handler, ok := gsa.configs.serverStream[req.Route]
	if !ok {
		return errors.New(fmt.Sprintf("Route not found %s.", req.Route))
	}

	// defer log
	defer gsa.log(stream.Context(), handler, startTime)

	errs := runtime.TaskSetHandler(stream.Context(), handler.TaskSet)
	if errs != nil {
		for _, err := range errs {
			handler.LogTaskError(stream.Context(), err)
		}
		return errors.New(fmt.Sprintf("Server streaming %s failed when handling tasks.", handler.GetId()))
	}

	payloadCount := handler.ServerConfig.Payload.Count
	if payloadCount <= 0 {
		payloadCount = constants.DefaultPayloadCount
	}

	for i := 0; i < payloadCount; i++ {
		// write response
		payloadSize := model.GetPayloadSize(handler.ServerConfig.Payload)
		payload := util.GenerateRandomString(payloadSize)

		rpcResponse := pb.StreamResponse{
			Payload: payload,
		}

		logger.Logger.Debugf("Ran %s, responding with %v bytes.", handler.GetId(), payloadSize)

		if err := stream.Send(&rpcResponse); err != nil {
			return err
		}
	}
	return nil
}

// Bidirectional streaming
func (gsa *GrpcServerAdapter) BidirectionalStreaming(stream pb.Grpc_BidirectionalStreamingServer) error {
	// record time for logging
	startTime := time.Now()

	// process the first message in the stream and start tasks
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	handler, ok := gsa.configs.biStream[req.Route]
	if !ok {
		return errors.New(fmt.Sprintf("Route not found %s.", req.Route))
	}

	// defer log
	defer gsa.log(stream.Context(), handler, startTime)

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
			break
		}
		if err != nil {
			return err
		}

		// write response
		payloadSize := model.GetPayloadSize(handler.ServerConfig.Payload)
		payload := util.GenerateRandomString(payloadSize)

		rpcResponse := pb.StreamResponse{
			Payload: payload,
		}

		logger.Logger.Debugf("Ran %s, responding with %v bytes.", handler.GetId(), payloadSize)

		if err := stream.Send(&rpcResponse); err != nil {
			return err
		}
	}
	return nil
}

func (gsa *GrpcServerAdapter) log(ctx context.Context, handler *ports.PrimaryHandler, startTime time.Time) {
	elapsed := time.Since(startTime).Milliseconds()
	unit := "ms"
	if elapsed == 0 {
		elapsed = time.Since(startTime).Microseconds()
		unit = "μs"
	}
	logger.Logger.WithContext(ctx).Infof("Invocation handled | %-40s | %10v %s", handler.GetId(), elapsed, unit)
}
