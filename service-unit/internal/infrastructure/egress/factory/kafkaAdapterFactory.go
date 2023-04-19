package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	kafkaAdapter "github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/producer_adapter/kafka"
)

func (egressAdapterDetails EgressAdapterDetails) kafkaEgressAdapterFactory() (core.EgressAdapter, error) {
	var kafkaEgressAdapter core.EgressAdapter
	var err error
	if kafkaWriterConnection, ok := egressAdapterDetails.connection.(*kafkaAdapter.KafkaWriterConnection); ok {
		kafkaEgressAdapter = kafkaAdapter.KafkaProducerAdapter{Writer: kafkaWriterConnection.Connection}
	} else {
		err = errors.New("Unmatched connection instance")
	}
	return kafkaEgressAdapter, err
}
