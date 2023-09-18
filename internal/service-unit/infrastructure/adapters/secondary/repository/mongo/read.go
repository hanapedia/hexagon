package mongo

import (
	"context"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/domain/contract"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/config"
	tracing "github.com/hanapedia/the-bench/internal/service-unit/infrastructure/telemetry/tracing/mongo"
	"github.com/hanapedia/the-bench/pkg/service-unit/utils"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReadAdapter struct {
	Name string
	Database string
	Client *mongo.Client
	Collection constants.RepositoryEntryVariant
	ports.SecondaryPortBase
}

// Read the document in the intial set of data
func (mra *MongoReadAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
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
	if err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: &foundRecord.Payload,
		Error: nil,
	}
}
