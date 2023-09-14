package mongo

import (
	"context"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure_new/adapters/secondary/config"
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
}

// Update or insert to random id in range from number of initial data to twice the size of the initial data
func (mra MongoWriteAdapter) Call(ctx context.Context) (string, error) {
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
		return "Failed to write an entry.", err
	}
	return "Successfully wrote an entry.", err
}

