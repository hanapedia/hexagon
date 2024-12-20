package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/core/runtime"
	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/internal/service-unit/domain/contract"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/server"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	util "github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// must implement ports.PrimaryPort
type RestServerAdapter struct {
	addr   string
	server *http.Server
	mux    *http.ServeMux
}

func NewRestServerAdapter() *RestServerAdapter {
	addr := config.GetRestServerAddr()
	mux := http.NewServeMux()
	server := &http.Server{Addr: addr, Handler: mux}

	adapter := RestServerAdapter{
		addr:   config.GetRestServerAddr(),
		server: server,
		mux:    mux,
	}

	return &adapter
}

func (rsa *RestServerAdapter) Serve(ctx context.Context, shutdownWg, readyWg *sync.WaitGroup) error {
	logger.Logger.Infof("Serving rest server at %s", rsa.addr)
	go func() {
		<-ctx.Done()
		logger.Logger.Infof("Context cancelled. Rest Server shutting down.")
		rsa.server.Shutdown(context.Background())
		shutdownWg.Done()
	}()

	listen, err := net.Listen("tcp", rsa.addr)
	if err != nil {
		return err
	}

	// Mark Ready
	readyWg.Done()

	err = rsa.server.Serve(listen)
	if err != nil && err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (rsa *RestServerAdapter) Register(handler *domain.PrimaryAdapterHandler) error {
	if handler.ServerConfig == nil {
		return errors.New(fmt.Sprintf("Invalid configuartion for handler %s.", handler.GetId()))
	}

	var handlerFunc http.Handler

	handlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.SetServerAdapterInProgress(domain.INC, handler.ServiceName, handler.ServerConfig)
		// record time for logging
		startTime := time.Now()
		// defer log
		defer server.Log(handler, startTime)

		// call tasks
		result := runtime.TaskSetHandler(r.Context(), handler)
		defer func() {
			// record metrics
			server.ObserveServerAdapterDuration(time.Since(startTime), handler.ServiceName, handler.ServerConfig, result.ShouldFail)
			server.SetServerAdapterInProgress(domain.DEC, handler.ServiceName, handler.ServerConfig)
		}()

		if result.ShouldFail {
			w.WriteHeader(http.StatusInternalServerError)
			restResponse := contract.RestResponseBody{}
			json.NewEncoder(w).Encode(restResponse)
			return
		}

		// write response
		payloadSize := model.GetPayloadSize(handler.ServerConfig.Payload)
		payload := util.GenerateRandomString(payloadSize)

		restResponse := contract.RestResponseBody{
			Payload: &payload,
		}

		logger.Logger.Debugf("Ran %s, responding with %v bytes.", handler.GetId(), payloadSize)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(restResponse)
	})

	if config.GetEnvs().TRACING {
		handlerFunc = otelhttp.NewHandler(handlerFunc, handler.GetId())
	}

	var err error
	switch handler.ServerConfig.Action {
	case constants.GET, constants.READ:
		rsa.mux.Handle(fmt.Sprintf("GET /%s", handler.ServerConfig.Route), handlerFunc)
	case constants.POST, constants.WRITE:
		rsa.mux.Handle(fmt.Sprintf("POST /%s", handler.ServerConfig.Route), handlerFunc)
	default:
		err = errors.New("Handler has no matching action")
	}
	return err
}
