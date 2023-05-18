package factory

import (
	"errors"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/producer_adapter/kafka"
)

func kafkaEgressAdapterFactory(adapterConfig model.BrokerEgressAdapterConfig, connection core.EgressConnection) (core.EgressAdapter, error) {
	var kafkaEgressAdapter core.EgressAdapter
	var err error
	if kafkaProducerConnection, ok := (connection).(kafka.KafkaProducerConnection); ok {
		kafkaEgressAdapter = kafka.KafkaProducerAdapter{Writer: kafkaProducerConnection.Connection}
	} else {
		err = errors.New("Unmatched connection instance")
	}
	return kafkaEgressAdapter, err
}
