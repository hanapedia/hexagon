package mongo

import (
	"context"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	context *context.Context
	Client  *mongo.Client
}

// Client client for mongo
func NewMongoClient(addr string) ports.EgressClient {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))
	if err != nil {
		logger.Logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return MongoClient{Client: client, context: &ctx}
}

func (mongoClient MongoClient) Close() {
	mongoClient.Client.Disconnect(*mongoClient.context)
}
