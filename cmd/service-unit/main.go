package main

import (
	"reflect"

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
	initialization.InitTracing(serviceUnitConfig.Name)

	serviceUnit := initialization.NewServiceUnit(serviceUnitConfig)
	logger.Logger.Println("Service unit loaded")

	// setup service unit
	serviceUnit.Setup()

	// create error channel and start service unit
	errChan := make(chan ports.PrimaryPortError)
	serviceUnit.Start(errChan)

	serverAdapterError := <-errChan
	logger.Logger.Fatalf("%s failed: %s", reflect.TypeOf(serverAdapterError.PrimaryPort).Elem().Name(), serverAdapterError.Error)
}
