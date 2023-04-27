package factory

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/producer_adapter/kafka"
)

func brokerEgressAdapterFactory(adapterConfig model.BrokerAdapterConfig, connection core.EgressConnection) (core.EgressAdapter, error) {
	switch adapterConfig.Variant {
	case constants.KAFKA:
		return kafkaEgressAdapterFactory(adapterConfig, connection)
	default:
		err := errors.New("No matching protocol found")
		return nil, err
	}

}

func upsertBrokerEgressConnection(adapterConfig model.BrokerAdapterConfig, connections *map[string]core.EgressConnection) core.EgressConnection {
	key := fmt.Sprintf("%s.%s", adapterConfig.Variant, adapterConfig.Topic)
	connection, ok := (*connections)[key]
	if ok {
		log.Printf("connection already exists reusing %v", reflect.TypeOf(connection))
		return connection
	}
	switch adapterConfig.Variant {
	case constants.KAFKA:
		kafkaConnection := kafka.NewKafkaConnection(constants.KafkaBrokerAddr, adapterConfig.Topic)
		log.Printf("created new connection %v", reflect.TypeOf(kafkaConnection))

		(*connections)[key] = kafkaConnection
		return kafkaConnection
	default:
		log.Fatalf("invalid protocol")
	}
	return connection
}
