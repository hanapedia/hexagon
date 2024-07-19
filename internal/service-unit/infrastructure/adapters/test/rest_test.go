package test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/primary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	restServer "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/server/rest"
	restClient "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/invocation/rest"
	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func TestRestServerAndClient(t *testing.T) {
	// 1. Setup server
	server := restServer.NewRestServerAdapter()
	server.Register(&domain.PrimaryHandler{
		ServiceName: "test",
		ServerConfig: &v1.ServerConfig{
			Variant: "rest",
			Action:  constants.GET,
			Route:   "get",
		},
		TaskSet: []domain.Task{},
	})
	server.Register(&domain.PrimaryHandler{
		ServiceName: "test",
		ServerConfig: &v1.ServerConfig{
			Variant: "rest",
			Action:  constants.POST,
			Route:   "post",
		},
		TaskSet: []domain.Task{},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	errChan := make(chan primary.PrimaryPortError)
	go func() {
		wg.Add(1)
		if err := server.Serve(ctx, &wg); err != nil {
			errChan <- primary.PrimaryPortError{PrimaryPort: server, Error: err}
		}
	}()

	// 2. Setup client
	client := restClient.NewRestClient()
	readAdapter, err := restClient.RestInvocationAdapterFactory(
		&v1.InvocationConfig{
			Variant: "rest",
			Service: "localhost",
			Action:  constants.GET,
			Route:   "get",
		},
		client,
	)
	if err != nil {
		t.Fail()
		logger.Logger.Error(err)
		return
	}

	writeAdapter, err := restClient.RestInvocationAdapterFactory(
		&v1.InvocationConfig{
			Variant: "rest",
			Service: "localhost",
			Action:  constants.POST,
			Route:   "post",
		},
		client,
	)
	if err != nil {
		t.Fail()
		logger.Logger.Error(err)
		return
	}

	// TODO: replace with healthcheck probe
	time.Sleep(2 * time.Second)
	res := readAdapter.Call(context.Background())
	if res.Error != nil {
		t.Fail()
		logger.Logger.Error(res.Error)
		return
	}
	res = writeAdapter.Call(context.Background())
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
