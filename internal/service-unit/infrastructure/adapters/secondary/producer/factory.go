package producer

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/producer/kafka"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func NewSecondaryAdapter(adapterConfig *model.ProducerConfig, client secondary.SecondaryAdapterClient) (secondary.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.KAFKA:
		return kafka.KafkaProducerAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found when creating producer adapter.")
		return nil, err
	}

}

func NewClient(adapterConfig *model.ProducerConfig) secondary.SecondaryAdapterClient {
	switch adapterConfig.Variant {
	case constants.KAFKA:
		kafkaClient := kafka.NewKafkaClient(config.GetKafkaBrokerAddr(), adapterConfig.Topic)
		return kafkaClient
	default:
		logger.Logger.Fatalf("invalid protocol")
		return nil
	}
}
