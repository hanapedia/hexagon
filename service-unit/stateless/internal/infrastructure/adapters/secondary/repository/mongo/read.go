package mongo

import (
	"context"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/contract"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/config"
	tracing "github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/telemetry/tracing/mongo"
	"github.com/hanapedia/the-bench/service-unit/stateless/pkg/utils"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReadAdapter struct {
	Name string
	Database string
	Client *mongo.Client
	Collection constants.RepositoryEntryVariant
}

// Read the document in the intial set of data
func (mra MongoReadAdapter) Call(ctx context.Context) (string, error) {
	// create span if tracing is enabled
	if config.GetEnvs().TRACING {
		span := tracing.CreateWriteSpan(ctx, mra.Name, mra.Database, string(mra.Collection))
		defer span.End()
	}

	db := mra.Client.Database("mongo")
	collection := db.Collection(string(mra.Collection))
	// Find a document
	var foundRecord contract.MongoRecord
	id := utils.GetRandomId(1, constants.NumInitialEntries)
	filter := bson.D{{Key: "id", Value: id}}

	err := collection.FindOne(ctx, filter).Decode(&foundRecord)
	return foundRecord.Payload, err
}
