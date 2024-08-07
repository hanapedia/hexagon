package kafka

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func KafkaProducerAdapterFactory(adapterConfig *model.ProducerConfig, client secondary.SecondaryAdapterClient) (secondary.SecodaryPort, error) {
	var kafkaAdapter secondary.SecodaryPort
	var err error
	if kafkaProducerClient, ok := (client).(*KafkaProducerClient); ok {
		kafkaAdapter = &kafkaProducerAdapter{
			writer:      kafkaProducerClient.Client,
			payloadSize: model.GetPayloadSize(adapterConfig.Payload),
		}
	} else {
		err = errors.New("Unmatched client instance")
	}

	// set destionation id
	kafkaAdapter.SetDestId(adapterConfig.GetId())

	logger.Logger.Debugf("Initialized kafka producer adapter: %s", adapterConfig.GetId())
	return kafkaAdapter, err
}
