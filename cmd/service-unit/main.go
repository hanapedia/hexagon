package main

import (
	"reflect"

	"github.com/hanapedia/the-bench/internal/service-unit/application/core/initialization"
	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
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
