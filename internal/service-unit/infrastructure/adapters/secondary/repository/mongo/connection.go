package mongo

import (
	"context"

	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	context *context.Context
	Client  *mongo.Client
}

// Client client for mongo
func NewMongoClient(addr string) *MongoClient {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))
	if err != nil {
		logger.Logger.Fatalf("Failed to connect to MongoDB. err=%v, addr=%s", err, addr)
	}
	mongoClient := MongoClient{Client: client, context: &ctx}
	return &mongoClient
}

func (mongoClient *MongoClient) Close() {
	mongoClient.Client.Disconnect(*mongoClient.context)
}
