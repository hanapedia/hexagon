package primary

import (
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/consumer/kafka"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/server/grpc"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/server/rest"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func NewServerAdapter(config *model.ServerConfig) ports.PrimaryPort {
	var serverAdapter ports.PrimaryPort

	switch config.Variant {
	case constants.REST:
		serverAdapter = rest.NewRestServerAdapter()
	case constants.GRPC:
		serverAdapter = grpc.NewGrpcServerAdapter()
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
