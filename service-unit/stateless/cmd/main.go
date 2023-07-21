package main

import (
	"reflect"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/app/usecases"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

func main() {
	serviceUnitConfig := usecases.GetConfig("yaml")

	// init telemetry
	usecases.InitTelemetry(serviceUnitConfig.Name)

	serviceUnit := usecases.NewServiceUnit(serviceUnitConfig)
	logger.Logger.Println("Service unit successfully loaded.")

	serviceUnit.Setup()

	errChan := make(chan core.IngressAdapterError)
	serviceUnit.Start(errChan)

	serverAdapterError := <-errChan
	logger.Logger.Fatalf("%s failed: %s", reflect.TypeOf(serverAdapterError.IngressAdapter).Elem().Name(), serverAdapterError.Error)
}
