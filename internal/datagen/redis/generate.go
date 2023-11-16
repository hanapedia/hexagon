package redis

import (
	"fmt"
	"os"

	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
)

type RedisData struct {
	key  string
	data string
}

// count is the number of entries to generate
// size is the size of each data for each entry
func GenerateRedisData(count int, sizeVariant constants.PayloadSizeVariant) []RedisData {
	dataSlice := make([]RedisData, count)
	size := constants.PayloadSizeMap[sizeVariant]
	for i := 1; i <= count; i++ {
		payload := utils.GenerateRandomString(size)
		dataSlice = append(dataSlice, RedisData{
			key:  fmt.Sprintf("%s%v", sizeVariant, i),
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
