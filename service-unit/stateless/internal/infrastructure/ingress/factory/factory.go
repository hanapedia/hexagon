package factory

import (
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/consumer_adapter/kafka"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/ingress/server_adapter/rest"
)

func NewServerAdapter(serverAdapterProtocol constants.StatelessAdapterVariant) *ports.PrimaryPort {
	var serverAdapter ports.PrimaryPort

	switch serverAdapterProtocol {
	case constants.REST:
		serverAdapter = rest.NewRestServerAdapter()
	default:
		logger.Logger.Fatal("Adapter currently unsupported.")
	}

	return &serverAdapter
}

func NewConsumerAdapter(protocol constants.BrokerAdapterVariant, action string) *ports.PrimaryPort {
	var consumerAdapter ports.PrimaryPort

	switch protocol {
	case constants.KAFKA:
		consumerAdapter = kafka.NewKafkaConsumerAdapter(action)
	default:
		logger.Logger.Fatal("Adapter currently unsupported.")
	}

	return &consumerAdapter
}

// Takes the pointer to the slice of ServerAdapters
// Update or insert ServiceAdapter based on the handler input.
// Does not return any value
func RegiserHandlerToIngressAdapter(serviceName string, serverAdapter *ports.PrimaryPort, handler *ports.PrimaryAdapter) error {
	err := (*serverAdapter).Register(serviceName, handler)

	return err
}
