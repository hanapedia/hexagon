package main

import (
	"reflect"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/app/usecases"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
)

func main() {
	serviceUnitConfig := usecases.GetConfig("yaml")
	serviceUnit := usecases.NewServiceUnit(serviceUnitConfig)
	logger.Logger.Println("Service unit successfully loaded.")

	serviceUnit.Setup()

	errChan := make(chan core.IngressAdapterError)
	serviceUnit.Start(errChan)

	serverAdapterError := <-errChan
	logger.Logger.Fatalf("%s failed: %s", reflect.TypeOf(serverAdapterError.IngressAdapter).Elem().Name(), serverAdapterError.Error)
}
