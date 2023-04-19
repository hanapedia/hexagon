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
	log.Println("Service unit successfully loaded.")

	errChan := make(chan core.IngressAdapterError)
	for protocol, serverAdapter := range serviceUnit.ServerAdapters {
		serverAdapterCopy := serverAdapter
		log.Printf("Serving '%s' server.", protocol)
		go func() {
			if err := (*serverAdapterCopy).Serve(); err != nil {
				errChan <- core.IngressAdapterError{IngressAdapter: serverAdapterCopy, Error: err}
			}
		}()
	}

	serverAdapterError := <-errChan
	log.Fatalf("%s failed: %s", reflect.TypeOf(serverAdapterError.IngressAdapter).Elem().Name(), serverAdapterError.Error)
}
