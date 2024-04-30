package test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	grpcServer "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/server/grpc"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	grpcClient "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/invocation/grpc"
	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func TestGrpcServerAndClient(t *testing.T) {
	// 1. Setup server
	server := grpcServer.NewGrpcServerAdapter()
	server.Register(&ports.PrimaryHandler{
		ServiceName: "test",
		ServerConfig: &v1.ServerConfig{
			Variant: "grpc",
			Action:  constants.SIMPLE_RPC,
			Route:   "simple",
		},
		TaskSet: []ports.Task{},
	})
	server.Register(&ports.PrimaryHandler{
		ServiceName: "test",
		ServerConfig: &v1.ServerConfig{
			Variant: "grpc",
			Action:  constants.BI_STREAM,
			Route:   "bistream",
		},
		TaskSet: []ports.Task{},
	})
	server.Register(&ports.PrimaryHandler{
		ServiceName: "test",
		ServerConfig: &v1.ServerConfig{
			Variant: "grpc",
			Action:  constants.CLIENT_STREAM,
			Route:   "cstream",
		},
		TaskSet: []ports.Task{},
	})
	server.Register(&ports.PrimaryHandler{
		ServiceName: "test",
		ServerConfig: &v1.ServerConfig{
			Variant: "grpc",
			Action:  constants.SERVER_STREAM,
			Route:   "sstream",
		},
		TaskSet: []ports.Task{},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	errChan := make(chan ports.PrimaryPortError)
	go func() {
		wg.Add(1)
		if err := server.Serve(ctx, &wg); err != nil {
			errChan <- ports.PrimaryPortError{PrimaryPort: server, Error: err}
		}
	}()

	// 2. Setup client
	simpleRpcInvocationConfig := &v1.InvocationConfig{
		Variant: "grpc",
		Service: "localhost",
		Action:  constants.SIMPLE_RPC,
		Route:   "simple",
	}

	biStreamInvocationConfig := &v1.InvocationConfig{
		Variant: "grpc",
		Service: "localhost",
		Action:  constants.BI_STREAM,
		Route:   "bistream",
	}

	cStreamInvocationConfig := &v1.InvocationConfig{
		Variant: "grpc",
		Service: "localhost",
		Action:  constants.CLIENT_STREAM,
		Route:   "cstream",
	}

	sStreamInvocationConfig := &v1.InvocationConfig{
		Variant: "grpc",
		Service: "localhost",
		Action:  constants.SERVER_STREAM,
		Route:   "sstream",
	}

	addr := config.GetGrpcDialAddr(simpleRpcInvocationConfig)
	client := grpcClient.NewGrpcClient(addr)

	simpleRpcAdapter, err := grpcClient.GrpcInvocationAdapterFactory(
		simpleRpcInvocationConfig,
		client,
	)
	if err != nil {
		t.Fail()
		logger.Logger.Error(err)
		return
	}

	biStreamAdapter, err := grpcClient.GrpcInvocationAdapterFactory(
		biStreamInvocationConfig,
		client,
	)
	if err != nil {
		t.Fail()
		logger.Logger.Error(err)
		return
	}

	cStreamAdapter, err := grpcClient.GrpcInvocationAdapterFactory(
		cStreamInvocationConfig,
		client,
	)
	if err != nil {
		t.Fail()
		logger.Logger.Error(err)
		return
	}

	sStreamAdapter, err := grpcClient.GrpcInvocationAdapterFactory(
		sStreamInvocationConfig,
		client,
	)
	if err != nil {
		t.Fail()
		logger.Logger.Error(err)
		return
	}

	// TODO: replace with healthcheck probe
	time.Sleep(2 * time.Second)

	res := simpleRpcAdapter.Call(context.Background())
	if res.Error != nil {
		t.Fail()
		logger.Logger.Error(res.Error)
		return
	}

	res = biStreamAdapter.Call(context.Background())
	if res.Error != nil {
		t.Fail()
		logger.Logger.Error(res.Error)
		return
	}

	res = cStreamAdapter.Call(context.Background())
	if res.Error != nil {
		t.Fail()
		logger.Logger.Error(res.Error)
		return
	}

	res = sStreamAdapter.Call(context.Background())
	if res.Error != nil {
		t.Fail()
		logger.Logger.Error(res.Error)
		return
	}

	go func() {
		serverAdapterError := <-errChan
		if serverAdapterError.Error != nil {
			t.Fail()
			logger.Logger.Error(serverAdapterError.Error)
			return
		}
	}()

	// 3. shutdown server
	logger.Logger.Info("Request successful. Cancelling server context for clean up.")
	cancel()
	wg.Wait()
	client.Close()
}

