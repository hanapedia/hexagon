package factory

import (
	"errors"
	"fmt"
	"reflect"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/producer_adapter/kafka"
)

func brokerEgressAdapterFactory(adapterConfig model.BrokerEgressAdapterConfig, connection core.EgressConnection) (core.EgressAdapter, error) {
	switch adapterConfig.Variant {
	case constants.KAFKA:
		return kafkaEgressAdapterFactory(adapterConfig, connection)
	default:
		err := errors.New("No matching protocol found when creating broker egress adapter.")
		return nil, err
	}

}

func upsertBrokerEgressConnection(adapterConfig model.BrokerEgressAdapterConfig, connections *map[string]core.EgressConnection) core.EgressConnection {
	key := fmt.Sprintf("%s.%s", adapterConfig.Variant, adapterConfig.Topic)
	connection, ok := (*connections)[key]
	if ok {
		logger.Logger.Infof("connection already exists reusing %v", reflect.TypeOf(connection))
		return connection
	}
	switch adapterConfig.Variant {
	case constants.KAFKA:
		kafkaConnection := kafka.NewKafkaConnection(config.GetKafkaBrokerAddr(), adapterConfig.Topic)
		logger.Logger.Infof("created new connection %v", reflect.TypeOf(kafkaConnection))

		(*connections)[key] = kafkaConnection
		return kafkaConnection
	default:
		logger.Logger.Fatalf("invalid protocol")
	}
	return connection
}
