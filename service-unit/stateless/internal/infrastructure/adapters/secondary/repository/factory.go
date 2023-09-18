package repository

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/repository/mongo"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

func NewSecondaryAdapter(adapterConfig *model.RepositoryClientConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.MONGO:
		return mongo.MongoClientAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found when creating repository client adapter.")
		return nil, err
	}

}

// NewClient creates new client to stateful service
func NewClient(adapterConfig *model.RepositoryClientConfig) ports.SecondaryAdapter {
	switch adapterConfig.Variant {
	case constants.MONGO:
		connectionUri := config.GetMongoConnectionUri(adapterConfig)
		mongoClient := mongo.NewMongoClient(connectionUri)
		return mongoClient
	default:
		logger.Logger.Fatalf("invalid protocol")
		return nil
	}
}
