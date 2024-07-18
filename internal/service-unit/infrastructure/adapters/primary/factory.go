package primary

import (
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/primary"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/consumer/kafka"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/server/grpc"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/server/rest"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func NewServerAdapter(config *model.ServerConfig) primary.PrimaryPort {
	var serverAdapter primary.PrimaryPort

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

func NewConsumerAdapter(config *model.ConsumerConfig, name string) primary.PrimaryPort {
	var consumerAdapter primary.PrimaryPort

	switch config.Variant {
	case constants.KAFKA:
		consumerAdapter = kafka.NewKafkaConsumerAdapter(config.Topic, config.GetConsumerGroupId(name))
	default:
		logger.Logger.Fatal("Adapter currently unsupported.")
	}

	return consumerAdapter
}
