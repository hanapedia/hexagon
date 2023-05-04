package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/egress/repository_adapter/mongo"
)

func mongoEgressAdapterFactory(adapterConfig model.StatefulAdapterConfig, connection core.EgressConnection) (core.EgressAdapter, error) {
	var mongoEgressAdapter core.EgressAdapter
	var err error
	if mongoConnection, ok := (connection).(mongo.MongoConnection); ok {
		switch adapterConfig.Action {
		case constants.READ:
			mongoEgressAdapter = mongo.MongoReadAdapter{
				Connection: mongoConnection.Connection,
				Collection: constants.RepositoryEntryVariant(adapterConfig.Size),
			}
		case constants.WRITE:
			mongoEgressAdapter = mongo.MongoWriteAdapter{
				Connection: mongoConnection.Connection,
				Collection: constants.RepositoryEntryVariant(adapterConfig.Size),
			}
		default:
			err = errors.New("No matching action found when creating mongo egress adapter.")
		}
	} else {
		err = errors.New("Unmatched connection instance")
	}
	return mongoEgressAdapter, err
}
