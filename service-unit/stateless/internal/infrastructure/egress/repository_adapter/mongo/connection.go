package mongo

import (
	"context"
	"log"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	context    *context.Context
	Connection *mongo.Client
}

// Client connection for mongo
func NewMongoConnection(addr string) core.EgressConnection {
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
    return MongoConnection{Connection: client, context: &ctx}
}

func (mongoConnection MongoConnection) Close() {
	mongoConnection.Connection.Disconnect(*mongoConnection.context)
}
