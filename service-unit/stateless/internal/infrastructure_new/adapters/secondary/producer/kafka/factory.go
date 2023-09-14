package kafka

import (
	"errors"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
)

func KafkaEgressAdapterFactory(adapterConfig model.BrokerEgressAdapterConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	var kafkaEgressAdapter ports.SecodaryPort
	var err error
	if kafkaProducerClient, ok := (client).(KafkaProducerClient); ok {
		kafkaEgressAdapter = KafkaProducerAdapter{Writer: kafkaProducerClient.Client}
	} else {
		err = errors.New("Unmatched client instance")
	}
	return kafkaEgressAdapter, err
}
