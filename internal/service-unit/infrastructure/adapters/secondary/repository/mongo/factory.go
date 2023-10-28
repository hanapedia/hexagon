package mongo

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
)

func MongoClientAdapterFactory(adapterConfig *model.RepositoryClientConfig, client ports.SecondaryAdapterClient) (ports.SecodaryPort, error) {
	var mongoAdapter ports.SecodaryPort
	var err error
	if mongoClient, ok := (client).(*MongoClient); ok {
		switch adapterConfig.Action {
		case constants.READ:
			mongoAdapter = &MongoReadAdapter{
				name:       adapterConfig.Name,
				database:   string(adapterConfig.Variant),
				client:     mongoClient.Client,
				collection: adapterConfig.Payload,
			}
		case constants.WRITE:
			mongoAdapter = &MongoWriteAdapter{
				name:       adapterConfig.Name,
				database:   string(adapterConfig.Variant),
				client:     mongoClient.Client,
				collection: adapterConfig.Payload,
			}
		default:
			err = errors.New("No matching action found when creating mongo client adapter.")
		}
	} else {
		err = errors.New("Unmatched client instance")
	}

	// set destionation id
	mongoAdapter.SetDestId(adapterConfig.GetId())

	return mongoAdapter, err
}
