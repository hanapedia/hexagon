package factory

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/producer_adapter/kafka"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/repository_adapter/mongo"
	"github.com/hanapedia/the-bench/config/constants"
)

type EgressAdapterDetails struct {
	serviceName string
	protocol    constants.AdapterProtocol
	action      string
	adapterName string
	connection  core.EgressConnection
}

func newEgressAdapterDetails(id string) (EgressAdapterDetails, error) {
	idSubstring := strings.Split(id, ".")
	var err error
	if len(idSubstring) != 4 {
		err = errors.New("Invalid adapter id")
	}
	return EgressAdapterDetails{
		serviceName: idSubstring[constants.ServiceNameIndex],
		protocol:    constants.AdapterProtocol(idSubstring[constants.ProtocolIndex]),
		action:      idSubstring[constants.ActionIndex],
		adapterName: idSubstring[constants.AdapterNameIndex],
	}, err
}

func NewEgressAdapter(id string, connections *map[string]core.EgressConnection) (core.EgressAdapter, error) {
	egressAdapterDetails, err := newEgressAdapterDetails(id)
	var egressAdapter core.EgressAdapter
	switch egressAdapterDetails.protocol {
	case constants.REST:
		egressAdapter, err = egressAdapterDetails.restEgressAdapterFactory()
	case constants.KAFKA:
		key := fmt.Sprintf("%s.%s", egressAdapterDetails.protocol, egressAdapterDetails.action)
		egressAdapterDetails.UpsertConnection(connections, key)
		egressAdapter, err = egressAdapterDetails.kafkaEgressAdapterFactory()
	case constants.MONGO:
		key := fmt.Sprintf("%s.%s", egressAdapterDetails.serviceName, egressAdapterDetails.protocol)
		egressAdapterDetails.UpsertConnection(connections, key)
		egressAdapter, err = egressAdapterDetails.mongoEgressAdapterFactory()
	default:
		err = errors.New("No matching protocol found")
	}

	return egressAdapter, err
}

func (egressAdapterDetails *EgressAdapterDetails) UpsertConnection(connections *map[string]core.EgressConnection, key string) {
	connection, ok := (*connections)[key]
	if ok {
		log.Printf("connection already exists reusing %v", reflect.TypeOf(connection))
		egressAdapterDetails.connection = connection
		return
	}
	switch egressAdapterDetails.protocol {
	case constants.KAFKA:
		kafkaConnection := kafka.NewKafkaConnection(constants.KafkaBrokerAddr, egressAdapterDetails.action)
		log.Printf("created new connection %v", reflect.TypeOf(kafkaConnection))

		egressAdapterDetails.connection = kafkaConnection
		(*connections)[key] = kafkaConnection
	case constants.MONGO:
		mongoConnection := mongo.NewMongoConnection(constants.MongoURIAddr)
		log.Printf("created new connection %v", reflect.TypeOf(mongoConnection))

		egressAdapterDetails.connection = mongoConnection
		(*connections)[key] = mongoConnection
	default:
		log.Fatalf("invalid protocol")
	}
}
