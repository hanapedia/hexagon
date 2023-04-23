package factory

import (
	"errors"
	"log"
	"reflect"

	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/repository_adapter/mongo"
)

func statefulEgressAdapterFactory(adapterConfig model.StatefulEgressConfig, connection core.EgressConnection) (core.EgressAdapter, error) {
	switch adapterConfig.Variant {
	case constants.MONGO:
		return mongoEgressAdapterFactory(adapterConfig, connection)
	default:
		err := errors.New("No matching protocol found")
		return nil, err
	}

}

func upsertStatefulEgressConnection(adapterConfig model.StatefulEgressConfig, connections *map[string]core.EgressConnection) core.EgressConnection {
	key := string(adapterConfig.Variant)
	connection, ok := (*connections)[key]
	if ok {
		log.Printf("connection already exists reusing %v", reflect.TypeOf(connection))
		return connection
	}
	switch adapterConfig.Variant {
	case constants.MONGO:
		mongoConnection := mongo.NewMongoConnection(constants.MongoURIAddr)
		log.Printf("created new connection %v", reflect.TypeOf(mongoConnection))

		(*connections)[key] = mongoConnection
		return mongoConnection
	default:
		log.Fatalf("invalid protocol")
	}
	return connection
}
