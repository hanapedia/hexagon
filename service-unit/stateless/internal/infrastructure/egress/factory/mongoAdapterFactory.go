package factory

import (
	"errors"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/repository_adapter/mongo"
)

func mongoEgressAdapterFactory(adapterConfig model.StatefulEgressAdapterConfig, client core.EgressClient) (core.EgressAdapter, error) {
	var mongoEgressAdapter core.EgressAdapter
	var err error
	if mongoClient, ok := (client).(mongo.MongoClient); ok {
		switch adapterConfig.Action {
		case constants.READ:
			mongoEgressAdapter = mongo.MongoReadAdapter{
				Client: mongoClient.Client,
				Collection: constants.RepositoryEntryVariant(adapterConfig.Size),
			}
		case constants.WRITE:
			mongoEgressAdapter = mongo.MongoWriteAdapter{
				Client: mongoClient.Client,
				Collection: constants.RepositoryEntryVariant(adapterConfig.Size),
			}
		default:
			err = errors.New("No matching action found when creating mongo egress adapter.")
		}
	} else {
		err = errors.New("Unmatched client instance")
	}
	return mongoEgressAdapter, err
}
