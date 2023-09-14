package factory

import (
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/server_adapter/rest"
)

func RestServerAdapterFactory(serviceName string, serverAdapters *[]*ports.PrimaryPort, handler *ports.PrimaryAdapter) {
	idx := -1
	for i, serverAdapter := range *serverAdapters {
		if restServerAdapter, ok := (*serverAdapter).(rest.RestServerAdapter); ok {
			restServerAdapter.Register(serviceName, handler)
			idx = i
			break
		}
	}

	if idx < 0 {
		restServerAdapter := rest.NewRestServerAdapter()
		restServerAdapter.Register(serviceName, handler)
	}
}
