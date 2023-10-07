package redis

import (
	"fmt"
	"os"

	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
	"github.com/hanapedia/the-bench/pkg/service-unit/payload"
)

type RedisData struct {
	key  string
	data string
}

// count is the number of entries to generate
// size is the size of each data for each entry
func GenerateRedisData(count int, size constants.PayloadSizeVariant) []RedisData {
	dataSlice := make([]RedisData, count)
	for i := 1; i <= count; i++ {
		payload, err := payload.GeneratePayload(size)
		if err != nil {
			logger.Logger.Panicf("Error generating random string %s", err)
		}
		dataSlice = append(dataSlice, RedisData{
			key:  fmt.Sprintf("%s%v", size, i),
			data: payload,
		})
	}
	return dataSlice
}

// WriteRedisDataToFile writes the key-value pairs of RedisData to a file.
func WriteRedisDataToFile(filename string, data []RedisData) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, entry := range data {
		_, err := file.WriteString(fmt.Sprintf("%s %s\n", entry.key, entry.data))
		if err != nil {
			return err
		}
	}

	return nil
}
