package primary

import (
	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/primary/consumer/kafka"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/primary/server/rest"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
)

func NewServerAdapter(config *model.ServerConfig) ports.PrimaryPort {
	var serverAdapter ports.PrimaryPort

	switch config.Variant {
	case constants.REST:
		serverAdapter = rest.NewRestServerAdapter()
	default:
		logger.Logger.Fatal("Adapter currently unsupported.")
	}

	return serverAdapter
}

func NewConsumerAdapter(config *model.ConsumerConfig) ports.PrimaryPort {
	var consumerAdapter ports.PrimaryPort

	switch config.Variant {
	case constants.KAFKA:
		consumerAdapter = kafka.NewKafkaConsumerAdapter(config.Topic)
	default:
		logger.Logger.Fatal("Adapter currently unsupported.")
	}

	return consumerAdapter
}

// Takes the pointer to the slice of ServerAdapters
// Update or insert ServiceAdapter based on the handler input.
// Does not return any value
func RegiserHandlerToPrimaryAdapter(serviceName string, serverAdapter ports.PrimaryPort, handler *ports.PrimaryHandler) error {
	err := serverAdapter.Register(handler)

	return err
}
