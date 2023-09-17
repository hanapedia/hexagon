package kafka

import (
	"errors"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
)

func KafkaProducerAdapterFactory(adapterConfig *model.ProducerConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	var kafkaAdapter ports.SecodaryPort
	var err error
	if kafkaProducerClient, ok := (client).(KafkaProducerClient); ok {
		kafkaAdapter = &KafkaProducerAdapter{
			Writer: kafkaProducerClient.Client,
		}
	} else {
		err = errors.New("Unmatched client instance")
	}

	// set destionation id
	kafkaAdapter.SetDestId(adapterConfig.GetId())

	return kafkaAdapter, err
}
