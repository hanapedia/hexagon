package redis

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func RedisClientAdapterFactory(adapterConfig *model.RepositoryClientConfig, client secondary.SecondaryAdapterClient) (secondary.SecodaryPort, error) {
	var redisAdapter secondary.SecodaryPort
	var err error
	if redisClient, ok := (client).(*RedisClient); ok {
		switch adapterConfig.Action {
		case constants.READ:
			redisAdapter = &redisReadAdapter{
				name:           adapterConfig.Name,
				client:         redisClient.Client,
				payloadVariant: model.GetPayloadVariant(adapterConfig.Payload),
			}
		case constants.WRITE:
			redisAdapter = &redisWriteAdapter{
				name:           adapterConfig.Name,
				client:         redisClient.Client,
				payloadSize:    model.GetPayloadSize(adapterConfig.Payload),
				payloadVariant: model.GetPayloadVariant(adapterConfig.Payload),
			}
		default:
			err = errors.New("No matching action found when creating redis client adapter.")
		}
	} else {
		err = errors.New("Unmatched client instance")
	}

	// set destionation id
	redisAdapter.SetDestId(adapterConfig.GetId())

	logger.Logger.Debugf("Initialized redis repository adapter: %s", adapterConfig.GetId())
	return redisAdapter, err
}
