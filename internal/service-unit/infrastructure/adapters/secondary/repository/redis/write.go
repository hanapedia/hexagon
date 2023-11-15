package redis

import (
	"context"
	"fmt"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"github.com/redis/go-redis/v9"
)

type redisWriteAdapter struct {
	name           string
	client         *redis.Client
	payloadSize    int64
	payloadVariant constants.PayloadSizeVariant
	ports.SecondaryPortBase
}

// Update or insert to random id in range from number of initial data to twice the size of the initial data
func (rwa *redisWriteAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	payload, err := utils.GenerateRandomString(rwa.payloadSize)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	id := utils.GetRandomId(constants.NumInitialEntries+1, constants.NumInitialEntries*2)
	key := fmt.Sprintf("%s%v", rwa.payloadVariant, id)
	err = rwa.client.Set(ctx, key, payload, 0).Err()
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: &payload,
		Error:   nil,
	}
}
