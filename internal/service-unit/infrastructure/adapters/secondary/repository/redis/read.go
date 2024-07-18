package redis

import (
	"context"
	"fmt"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
	"github.com/redis/go-redis/v9"
)

type redisReadAdapter struct {
	name   string
	client *redis.Client
	payloadVariant constants.PayloadSizeVariant
	secondary.SecondaryPortBase
}

// Read the document in the intial set of data
func (rra *redisReadAdapter) Call(ctx context.Context) secondary.SecondaryPortCallResult {
	id := utils.GetRandomId(1, constants.NumInitialEntries)
	key := fmt.Sprintf("%s%v", rra.payloadVariant, id)
	val, err := rra.client.Get(ctx, key).Result()
	if err != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return secondary.SecondaryPortCallResult{
		Payload: &val,
		Error:   nil,
	}
}
