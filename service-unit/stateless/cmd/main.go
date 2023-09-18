package main

import (
	"reflect"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/core/initialization"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

func main() {
	// load config from yaml
	yamlConfigLoader := config.NewConfigLoader("yaml")
	serviceUnitConfig := initialization.GetConfig(yamlConfigLoader)

	// init telemetry
	initialization.InitTelemetry(serviceUnitConfig.Name)

	serviceUnit := initialization.NewServiceUnit(serviceUnitConfig)
	logger.Logger.Println("Service unit successfully loaded.")

	// setup service unit
	serviceUnit.Setup()

	// create error channel and start service unit
	errChan := make(chan ports.PrimaryPortError)
	serviceUnit.Start(errChan)

	serverAdapterError := <-errChan
	logger.Logger.Fatalf("%s failed: %s", reflect.TypeOf(serverAdapterError.PrimaryPort).Elem().Name(), serverAdapterError.Error)
}
