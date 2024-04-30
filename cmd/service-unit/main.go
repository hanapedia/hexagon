package main

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/hanapedia/hexagon/internal/service-unit/application/core/initialization"
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func main() {
	// set log level
	initialization.InitLogging()

	// load config from yaml
	yamlConfigLoader := config.NewConfigLoader("yaml")
	serviceUnitConfig := initialization.GetConfig(yamlConfigLoader)

	// init telemetry
	traceProvider := initialization.InitTracing(serviceUnitConfig.Name)
	if traceProvider != nil {
		defer traceProvider.Shutdown(context.Background())
	}

	serviceUnit := initialization.NewServiceUnit(serviceUnitConfig)
	logger.Logger.Println("Service unit loaded")

	// setup service unit
	serviceUnit.Setup()

	// create wait group to wait for graceful shutdown
    var wg sync.WaitGroup

	// create context for triggering graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle OS signals
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigs
        logger.Logger.Info("Termination Signal received, cancelling context.")
        cancel()
    }()

	// crate error chan
	errChan := make(chan ports.PrimaryPortError)
	go func() {
		serverAdapterError := <-errChan
		logger.Logger.Fatalf("%s failed: %s", reflect.TypeOf(serverAdapterError.PrimaryPort).Elem().Name(), serverAdapterError.Error)
	}()

	// start primary adapters
	serviceUnit.Start(ctx, &wg, errChan)
	wg.Wait()

	// Close all client connections
	serviceUnit.Close()
	return
}
