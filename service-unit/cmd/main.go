package main

import (
	"log"
	"reflect"

	"github.com/hanapedia/the-bench/service-unit/internal/app/usecases"
	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
)

func main() {
	configLoader := usecases.NewConfigLoader("yaml")
	serviceUnit := usecases.NewServiceUnit(configLoader)

	errChan := make(chan core.ServerAdapterError)
	for _, serverAdapter := range *serviceUnit.ServerAdapters {
		serverAdapterCopy := serverAdapter
		go func() {
			if err := (*serverAdapterCopy).Serve(); err != nil {
				errChan <- core.ServerAdapterError{ServerAdapter: *serverAdapterCopy, Error: err}
			}
		}()
	}

	serverAdapterError := <-errChan
	log.Fatalf("%s failed: %s", reflect.TypeOf(serverAdapterError.ServerAdapter).Elem().Name(), serverAdapterError.Error)
}
