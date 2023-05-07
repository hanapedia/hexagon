package factory

import (
	"errors"
	"reflect"

	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/config/logger"
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/repository_adapter/mongo"
)

func statefulEgressAdapterFactory(adapterConfig model.StatefulEgressAdapterConfig, connection core.EgressConnection) (core.EgressAdapter, error) {
	switch adapterConfig.Variant {
	case constants.MONGO:
		return mongoEgressAdapterFactory(adapterConfig, connection)
	default:
		err := errors.New("No matching protocol found when creating stateful egress adapter.")
		return nil, err
	}

}

func upsertStatefulEgressConnection(adapterConfig model.StatefulEgressAdapterConfig, connections *map[string]core.EgressConnection) core.EgressConnection {
	key := string(adapterConfig.Variant)
	connection, ok := (*connections)[key]
	if ok {
		logger.Logger.Infof("connection already exists reusing %v", reflect.TypeOf(connection))
		return connection
	}
	switch adapterConfig.Variant {
	case constants.MONGO:
		connectionUri := config.GetMongoConnectionUri(adapterConfig)
		mongoConnection := mongo.NewMongoConnection(connectionUri)
		logger.Logger.Infof("created new connection %v", reflect.TypeOf(mongoConnection))

		(*connections)[key] = mongoConnection
		return mongoConnection
	default:
		logger.Logger.Fatalf("invalid protocol")
	}
	return connection
}

