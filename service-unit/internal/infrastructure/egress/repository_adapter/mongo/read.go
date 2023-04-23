package mongo

import (
	"context"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/contract"
	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/service-unit/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReadAdapter struct {
	Connection *mongo.Client
	Collection constants.RepositoryEntryVariant
}

// Read the document in the intial set of data
func (mra MongoReadAdapter) Call() (string, error) {
	db := mra.Connection.Database("mongo")
	collection := db.Collection(string(mra.Collection))
	// Find a document
	var foundRecord contract.MongoRecord
	id := utils.GetRandomId(1, constants.NumInitialEntries)
	filter := bson.D{{Key: "id", Value: id}}

	ctx := context.Background()

	err := collection.FindOne(ctx, filter).Decode(&foundRecord)
	return foundRecord.Payload, err
}
