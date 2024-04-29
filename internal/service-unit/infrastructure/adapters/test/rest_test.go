package test

import (
	"context"
	"testing"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	restServer "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/server/rest"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/invocation/rest"
	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func TestRestServerAndClient(t *testing.T) {
	server := restServer.NewRestServerAdapter()
	server.Register(&ports.PrimaryHandler{
		ServiceName: "test",
		ServerConfig: &v1.ServerConfig{
			Variant: "rest",
			Action: constants.GET,
			Route: "get",
		},
		TaskSet: []ports.Task{},
	})
	server.Register(&ports.PrimaryHandler{
		ServiceName: "test",
		ServerConfig: &v1.ServerConfig{
			Variant: "rest",
			Action: constants.POST,
			Route: "post",
		},
		TaskSet: []ports.Task{},
	})
	errChan := make(chan ports.PrimaryPortError)
	go func() {
		// TODO: handle graceful shut down
		if err := server.Serve(); err != nil {
			errChan <- ports.PrimaryPortError{PrimaryPort: server, Error: err}
		}
	}()

	client := rest.NewRestClient()
	readAdapter, err := rest.RestInvocationAdapterFactory(
		&v1.InvocationConfig{
			Variant: "rest",
			Service: "localhost",
			Action: constants.GET,
			Route: "get",
		},
		client,
	)
	if err != nil {
		t.Fail()
		logger.Logger.Error(err)
		return
	}

	res := readAdapter.Call(context.Background())
	if res.Error != nil {
		t.Fail()
		logger.Logger.Error(res.Error)
		return
	}

	serverAdapterError := <-errChan
	if serverAdapterError.Error != nil {
		t.Fail()
		logger.Logger.Error(serverAdapterError.Error)
		return
	}
}
