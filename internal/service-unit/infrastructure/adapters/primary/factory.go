package primary

import (
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/primary/consumer/kafka"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/primary/server/rest"
)

func NewServerAdapter(serverAdapterProtocol constants.SeverAdapterVariant) ports.PrimaryPort {
	var serverAdapter ports.PrimaryPort

	switch serverAdapterProtocol {
	case constants.REST:
		serverAdapter = rest.NewRestServerAdapter()
	default:
		logger.Logger.Fatal("Adapter currently unsupported.")
	}

	return serverAdapter
}

func NewConsumerAdapter(protocol constants.BrokerVariant, action string) ports.PrimaryPort {
	var consumerAdapter ports.PrimaryPort

	switch protocol {
	case constants.KAFKA:
		consumerAdapter = kafka.NewKafkaConsumerAdapter(action)
	default:
		logger.Logger.Fatal("Adapter currently unsupported.")
	}

	return consumerAdapter
}

// Takes the pointer to the slice of ServerAdapters
// Update or insert ServiceAdapter based on the handler input.
// Does not return any value
func RegiserHandlerToPrimaryAdapter(serviceName string, serverAdapter ports.PrimaryPort, handler *ports.PrimaryHandler) error {
	err := serverAdapter.Register(serviceName, handler)

	return err
}
