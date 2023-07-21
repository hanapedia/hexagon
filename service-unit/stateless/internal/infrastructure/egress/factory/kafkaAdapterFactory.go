package factory

import (
	"errors"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/producer_adapter/kafka"
)

func kafkaEgressAdapterFactory(adapterConfig model.BrokerEgressAdapterConfig, client core.EgressClient) (core.EgressAdapter, error) {
	var kafkaEgressAdapter core.EgressAdapter
	var err error
	if kafkaProducerClient, ok := (client).(kafka.KafkaProducerClient); ok {
		kafkaEgressAdapter = kafka.KafkaProducerAdapter{Writer: kafkaProducerClient.Client}
	} else {
		err = errors.New("Unmatched client instance")
	}
	return kafkaEgressAdapter, err
}
