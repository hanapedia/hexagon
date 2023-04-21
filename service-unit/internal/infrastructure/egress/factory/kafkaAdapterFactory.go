package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/producer_adapter/kafka"
)

func (egressAdapterDetails EgressAdapterDetails) kafkaEgressAdapterFactory() (core.EgressAdapter, error) {
	var kafkaEgressAdapter core.EgressAdapter
	var err error
	if kafkaProducerConnection, ok := (egressAdapterDetails.connection).(kafka.KafkaProducerConnection); ok {
		kafkaEgressAdapter = kafka.KafkaProducerAdapter{Writer: kafkaProducerConnection.Connection}
	} else {
		err = errors.New("Unmatched connection instance")
	}
	return kafkaEgressAdapter, err
}
