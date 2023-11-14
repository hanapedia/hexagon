package mongo

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	"github.com/hanapedia/hexagon/internal/service-unit/domain/contract"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	tracing "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/telemetry/tracing/mongo"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoReadAdapter struct {
	name       string
	database   string
	client     *mongo.Client
	collection constants.PayloadSizeVariant
	ports.SecondaryPortBase
}

// Read the document in the intial set of data
func (mra *mongoReadAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	// create span if tracing is enabled
	if config.GetEnvs().TRACING {
		span := tracing.CreateWriteSpan(ctx, mra.name, mra.database, string(mra.collection))
		defer span.End()
	}

	db := mra.client.Database("mongo")
	collection := db.Collection(string(mra.collection))
	// Find a document
	var foundRecord contract.MongoRecord
	id := utils.GetRandomId(1, constants.NumInitialEntries)
	filter := bson.D{{Key: "id", Value: id}}

	err := collection.FindOne(ctx, filter).Decode(&foundRecord)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: &foundRecord.Payload,
		Error:   nil,
	}
}
