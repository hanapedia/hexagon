package mongo

import (
	"context"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/adapters/secondary/config"
	tracing "github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/telemetry/tracing/mongo"
	"github.com/hanapedia/the-bench/service-unit/stateless/pkg/utils"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoWriteAdapter struct {
	Name string
	Database string
	Client *mongo.Client
	Collection constants.RepositoryEntryVariant
	ports.SecondaryPortBase
}

// Update or insert to random id in range from number of initial data to twice the size of the initial data
func (mra *MongoWriteAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	// create span if tracing is enabled
	if config.GetEnvs().TRACING {
		span := tracing.CreateWriteSpan(ctx, mra.Name, mra.Database, string(mra.Collection))
		defer span.End()
	}

	db := mra.Client.Database(mra.Database)
	collection := db.Collection(string(mra.Collection))

	payload, err := utils.GeneratePayloadWithRepositorySize(mra.Collection)
	id := utils.GetRandomId(constants.NumInitialEntries+1, constants.NumInitialEntries*2)
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"payload": payload}}
	updateOpts := options.Update().SetUpsert(true)

	_, err = collection.UpdateOne(ctx, filter, update, updateOpts)
	if err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: &payload,
		Error: nil,
	}
}

