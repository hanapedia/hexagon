package factory

import (
	"log"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/consumer_adapter/kafka"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/server_adapter/rest"
	"github.com/hanapedia/the-bench/config/constants"
)

func NewServerAdapter(serverAdapterProtocol constants.StatelessAdapterVariant) *core.IngressAdapter {
	var serverAdapter core.IngressAdapter

	switch serverAdapterProtocol {
	case constants.REST:
		serverAdapter = rest.NewRestServerAdapter()
	default:
		log.Fatal("Adapter currently unsupported.")
	}

	return &serverAdapter
}

func NewConsumerAdapter(protocol constants.BrokerAdapterVariant, action string) *core.IngressAdapter {
	var consumerAdapter core.IngressAdapter

	switch protocol {
	case constants.KAFKA:
		consumerAdapter = kafka.NewKafkaConsumerAdapter(action)
	default:
		log.Fatal("Adapter currently unsupported.")
	}

	return &consumerAdapter
}

// Takes the pointer to the slice of ServerAdapters
// Update or insert ServiceAdapter based on the handler input.
// Does not return any value
func RegiserHandlerToIngressAdapter(serverAdapter *core.IngressAdapter, handler *core.IngressAdapterHandler) error {
	err := (*serverAdapter).Register(handler)

	return err
}
