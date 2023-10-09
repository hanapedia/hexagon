package redis

import (
	"context"
	"fmt"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/pkg/common/utils"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/service-unit/payload"
	"github.com/redis/go-redis/v9"
)

type redisWriteAdapter struct {
	name   string
	client *redis.Client
	size   constants.PayloadSizeVariant
	ports.SecondaryPortBase
}

// Update or insert to random id in range from number of initial data to twice the size of the initial data
func (rwa *redisWriteAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	payload, err := payload.GeneratePayload(rwa.size)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	id := utils.GetRandomId(constants.NumInitialEntries+1, constants.NumInitialEntries*2)
	key := fmt.Sprintf("%s%v", rwa.size, id)
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