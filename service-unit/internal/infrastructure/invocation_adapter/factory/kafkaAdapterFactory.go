package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	kafkaAdapter "github.com/hanapedia/the-bench/service-unit/internal/infrastructure/invocation_adapter/kafka"
)

func (invocationAdapterDetails InvocationAdapterDetails) kafkaInvocationAdapterFactory() (core.InvocationAdapter, error) {
	var kafkaInvocationAdapter core.InvocationAdapter
	var err error
	if kafkaWriterConnection, ok := invocationAdapterDetails.connection.(*kafkaAdapter.KafkaWriterConnection); ok {
		kafkaInvocationAdapter = kafkaAdapter.KafkaProducerAdapter{Writer: kafkaWriterConnection.Connection}
	} else {
		err = errors.New("Unmatched connection instance")
	}
	return kafkaInvocationAdapter, err
}
