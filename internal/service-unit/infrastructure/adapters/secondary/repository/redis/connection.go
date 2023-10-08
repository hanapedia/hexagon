package redis

import (
	"context"

	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	context *context.Context
	Client  *redis.Client
}

// Client client for redis
func NewRedisClient(addr string) *redisClient {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	// propagate trace header if tracing is enabled
	if config.GetEnvs().TRACING {
	// Enable tracing instrumentation.
		if err := redisotel.InstrumentTracing(client); err != nil {
			logger.Logger.Panicf("Failed to enable tracing for redis. addr=%s, err=%s", addr, err)
		}
	}

	redisClient := redisClient{Client: client, context: &ctx}
	return &redisClient
}

func (redisClient *redisClient) Close() {
	redisClient.Client.Close()
}
