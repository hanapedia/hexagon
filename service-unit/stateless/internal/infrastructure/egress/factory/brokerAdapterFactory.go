package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure_new/adapters/secondary/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/producer_adapter/kafka"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

func brokerEgressAdapterFactory(adapterConfig model.BrokerEgressAdapterConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.KAFKA:
		return kafkaEgressAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found when creating broker egress adapter.")
		return nil, err
	}

}

func getOrCreateBrokerEgressClient(adapterConfig model.BrokerEgressAdapterConfig, clients *map[string]ports.SecondaryAdapter) ports.SecondaryAdapter {
	key := adapterConfig.GetId()
	client, ok := (*clients)[key]
	if ok {
		logger.Logger.Infof("client already exists reusing %v", key)
		return client
	}
	switch adapterConfig.Variant {
	case constants.KAFKA:
		kafkaClient := kafka.NewKafkaClient(config.GetKafkaBrokerAddr(), adapterConfig.Topic)
		logger.Logger.Infof("created new client %s", key)

		(*clients)[key] = kafkaClient
		return kafkaClient
	default:
		logger.Logger.Fatalf("invalid protocol")
	}
	return client
}
