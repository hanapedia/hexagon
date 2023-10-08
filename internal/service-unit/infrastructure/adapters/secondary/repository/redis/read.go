package redis

import (
	"context"
	"fmt"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/pkg/common/utils"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/redis/go-redis/v9"
)

type redisReadAdapter struct {
	name       string
	client     *redis.Client
	size constants.PayloadSizeVariant
	ports.SecondaryPortBase
}

// Read the document in the intial set of data
func (rra *redisReadAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	id := utils.GetRandomId(1, constants.NumInitialEntries)
	key := fmt.Sprintf("%s%v", rra.size, id)
    val, err := rra.client.Get(ctx, key).Result()
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: &val,
		Error:   nil,
	}
}
