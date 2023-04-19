package factory

import (
	"log"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/ingress/consumer_adapter/kafka"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/ingress/server_adapter/rest"
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
)

func NewServerAdapter(serverAdapterProtocol constants.AdapterProtocol) *core.IngressAdapter {
	var serverAdapter core.IngressAdapter

	switch serverAdapterProtocol {
	case constants.REST:
		serverAdapter = rest.NewRestServerAdapter()
	case constants.KAFKA:
		serverAdapter = kafka.NewKafkaConsumerAdapter()
	default:
		log.Fatal("Adapter currently unsupported.")
	}

	return &serverAdapter
}

// Takes the pointer to the slice of ServerAdapters
// Update or insert ServiceAdapter based on the handler input.
// Does not return any value
func RegiserHandlerToServerAdapter(serverAdapter *core.IngressAdapter, handler *core.Handler) error {
	err := (*serverAdapter).Register(handler)

	return err
}
