package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/egress/repository_adapter/mongo"
	"github.com/hanapedia/the-bench/config/constants"
)

func (egressAdapterDetails EgressAdapterDetails) mongoEgressAdapterFactory() (core.EgressAdapter, error) {
	var mongoEgressAdapter core.EgressAdapter
	var err error
	if mongoConnection, ok := (egressAdapterDetails.connection).(mongo.MongoConnection); ok {
		switch egressAdapterDetails.action {
		case string(constants.READ):
			mongoEgressAdapter = mongo.MongoReadAdapter{
				Connection: mongoConnection.Connection,
				Collection: constants.RepositoryEntryVariant(egressAdapterDetails.adapterName),
			}
		case string(constants.WRITE):
			mongoEgressAdapter = mongo.MongoWriteAdapter{
				Connection: mongoConnection.Connection,
				Collection: constants.RepositoryEntryVariant(egressAdapterDetails.adapterName),
			}
		default:
			err = errors.New("No matching protocol found")
		}
	} else {
		err = errors.New("Unmatched connection instance")
	}
	return mongoEgressAdapter, err
}
