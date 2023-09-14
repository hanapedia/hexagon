package mongo

import (
	"errors"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
)

func MongoEgressAdapterFactory(adapterConfig model.StatefulEgressAdapterConfig, client ports.SecondaryAdapter) (ports.SecodaryPort, error) {
	var mongoEgressAdapter ports.SecodaryPort
	var err error
	if mongoClient, ok := (client).(MongoClient); ok {
		switch adapterConfig.Action {
		case constants.READ:
			mongoEgressAdapter = MongoReadAdapter{
				Name: adapterConfig.Name,
				Database: string(adapterConfig.Variant),
				Client: mongoClient.Client,
				Collection: constants.RepositoryEntryVariant(adapterConfig.Size),
			}
		case constants.WRITE:
			mongoEgressAdapter = MongoWriteAdapter{
				Name: adapterConfig.Name,
				Database: string(adapterConfig.Variant),
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
