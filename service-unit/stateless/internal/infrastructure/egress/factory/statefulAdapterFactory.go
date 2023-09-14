package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure_new/adapters/secondary/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/repository_adapter/mongo"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

func statefulEgressAdapterFactory(adapterConfig model.StatefulEgressAdapterConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.MONGO:
		return mongoEgressAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found when creating stateful egress adapter.")
		return nil, err
	}

}

// getOrCreateStatefulEgressClient creates new client to stateful service if it does not exist
func getOrCreateStatefulEgressClient(adapterConfig model.StatefulEgressAdapterConfig, clients *map[string]ports.SecondaryAdapter) ports.SecondaryAdapter {
	key := adapterConfig.GetId()
	client, ok := (*clients)[key]
	if ok {
		logger.Logger.Infof("Client already exists reusing %s", key)
		return client
	}
	switch adapterConfig.Variant {
	case constants.MONGO:
		connectionUri := config.GetMongoConnectionUri(adapterConfig)
		mongoClient := mongo.NewMongoClient(connectionUri)
		logger.Logger.Infof("created new client %s", key)

		(*clients)[key] = mongoClient
		return mongoClient
	default:
		logger.Logger.Fatalf("invalid protocol")
	}
	return client
}
