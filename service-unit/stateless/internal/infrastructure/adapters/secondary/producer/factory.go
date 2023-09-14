package producer

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/producer/kafka"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

func NewSecondaryAdapter(adapterConfig model.BrokerEgressAdapterConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.KAFKA:
		return kafka.KafkaEgressAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found when creating broker egress adapter.")
		return nil, err
	}

}

func GetOrCreateClient(adapterConfig model.BrokerEgressAdapterConfig, clients *map[string]ports.SecondaryAdapter) ports.SecondaryAdapter {
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
