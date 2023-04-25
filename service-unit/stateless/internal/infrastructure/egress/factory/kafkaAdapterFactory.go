package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/producer_adapter/kafka"
)

func kafkaEgressAdapterFactory(adapterConfig model.BrokerEgressConfig, connection core.EgressConnection) (core.EgressAdapter, error) {
	var kafkaEgressAdapter core.EgressAdapter
	var err error
	if kafkaProducerConnection, ok := (connection).(kafka.KafkaProducerConnection); ok {
		kafkaEgressAdapter = kafka.KafkaProducerAdapter{Writer: kafkaProducerConnection.Connection}
	} else {
		err = errors.New("Unmatched connection instance")
	}
	return kafkaEgressAdapter, err
}
