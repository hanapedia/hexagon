package mongo

import (
	"context"

	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/service-unit/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoWriteAdapter struct {
	Connection *mongo.Client
	Collection constants.RepositoryEntryVariant
}

// Update or insert to random id in range from number of initial data to twice the size of the initial data
func (mra MongoWriteAdapter) Call() (string, error) {
	db := mra.Connection.Database("mongo")
	collection := db.Collection(string(mra.Collection))

	payload, err := utils.GeneratePayloadWithRepositorySize(mra.Collection)
	id := utils.GetRandomId(constants.NumInitialEntries+1, constants.NumInitialEntries*2)
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"payload": payload}}
	updateOpts := options.Update().SetUpsert(true)

	ctx := context.Background()

	_, err = collection.UpdateOne(ctx, filter, update, updateOpts)
	if err != nil {
		return "Failed to write an entry.", err
	}
	return "Successfully wrote an entry.", err
}

