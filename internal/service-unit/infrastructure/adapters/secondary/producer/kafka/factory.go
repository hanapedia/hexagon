package kafka

import (
	"errors"

	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
)

func KafkaProducerAdapterFactory(adapterConfig *model.ProducerConfig, client ports.SecondaryAdapterClient) (ports.SecodaryPort, error) {
	var kafkaAdapter ports.SecodaryPort
	var err error
	if kafkaProducerClient, ok := (client).(*KafkaProducerClient); ok {
		kafkaAdapter = &KafkaProducerAdapter{
			writer: kafkaProducerClient.Client,
			payload: adapterConfig.Payload,
		}
	} else {
		err = errors.New("Unmatched client instance")
	}

	// set destionation id
	kafkaAdapter.SetDestId(adapterConfig.GetId())

	return kafkaAdapter, err
}
