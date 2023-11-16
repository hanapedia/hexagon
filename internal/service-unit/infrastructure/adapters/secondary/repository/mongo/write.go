package mongo

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	tracing "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing/mongo"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoWriteAdapter struct {
	name       string
	database   string
	client     *mongo.Client
	collection constants.PayloadSizeVariant
	ports.SecondaryPortBase
}

// Update or insert to random id in range from number of initial data to twice the size of the initial data
func (mra *mongoWriteAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	// create span if tracing is enabled
	if config.GetEnvs().TRACING {
		span := tracing.CreateWriteSpan(ctx, mra.name, mra.database, string(mra.collection))
		defer span.End()
	}

	db := mra.client.Database(mra.database)
	collection := db.Collection(string(mra.collection))

	payload := utils.GenerateRandomString(constants.PayloadSizeMap[mra.collection])
	id := utils.GetRandomId(constants.NumInitialEntries+1, constants.NumInitialEntries*2)
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"payload": payload}}
	updateOpts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, updateOpts)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: &payload,
		Error:   nil,
	}
}
