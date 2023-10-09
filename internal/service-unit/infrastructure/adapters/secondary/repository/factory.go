package repository

import (
	"errors"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/repository/mongo"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/repository/redis"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
)

func NewSecondaryAdapter(adapterConfig *model.RepositoryClientConfig, client ports.SecondaryAdapterClient) (ports.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.MONGO:
		return mongo.MongoClientAdapterFactory(adapterConfig, client)
	case constants.REDIS:
		return redis.RedisClientAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found when creating repository client adapter.")
		return nil, err
	}

}

// NewClient creates new client to stateful service
func NewClient(adapterConfig *model.RepositoryClientConfig) ports.SecondaryAdapterClient {
	switch adapterConfig.Variant {
	case constants.MONGO:
		connectionUri := config.GetMongoConnectionUri(adapterConfig)
		mongoClient := mongo.NewMongoClient(connectionUri)
		return mongoClient
	case constants.REDIS:
		connectionUri := config.GetRedisConnectionAddr(adapterConfig)
		redisClient := redis.NewRedisClient(connectionUri)
		return redisClient
	default:
		logger.Logger.Fatalf("invalid protocol")
		return nil
	}
}
