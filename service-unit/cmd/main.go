package main

import (
	"log"
	"reflect"

	"github.com/hanapedia/the-bench/service-unit/internal/app/usecases"
	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
)

func main() {
	serviceUnit := usecases.NewServiceUnit("yaml")
	log.Println("Service unit successfully loaded.")

    serviceUnit.Setup()

	errChan := make(chan core.IngressAdapterError)
    serviceUnit.Start(errChan)

	serverAdapterError := <-errChan
	log.Fatalf("%s failed: %s", reflect.TypeOf(serverAdapterError.IngressAdapter).Elem().Name(), serverAdapterError.Error)
}
