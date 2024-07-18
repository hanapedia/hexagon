package redis

import (
	"context"
	"fmt"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"github.com/redis/go-redis/v9"
)

type redisWriteAdapter struct {
	name           string
	client         *redis.Client
	payloadSize    int64
	payloadVariant constants.PayloadSizeVariant
	secondary.SecondaryPortBase
}

// Update or insert to random id in range from number of initial data to twice the size of the initial data
func (rwa *redisWriteAdapter) Call(ctx context.Context) secondary.SecondaryPortCallResult {
	payload := utils.GenerateRandomString(rwa.payloadSize)

	id := utils.GetRandomId(constants.NumInitialEntries+1, constants.NumInitialEntries*2)
	key := fmt.Sprintf("%s%v", rwa.payloadVariant, id)
	err := rwa.client.Set(ctx, key, payload, 0).Err()
	if err != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return secondary.SecondaryPortCallResult{
		Payload: &payload,
		Error:   nil,
	}
}
