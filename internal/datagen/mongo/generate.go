package mongo

import (
	"encoding/json"
	"os"

	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/operator/logger"
	"github.com/hanapedia/the-bench/pkg/service-unit/payload"
)

type MongoData struct {
	ID      int    `json:"id"`
	Payload string `json:"payload"`
}

// count is the number of entries to generate
// size is the size of each data for each entry
func GenerateMongoData(count int, size constants.PayloadSizeVariant) []MongoData {
	dataSlice := make([]MongoData, count)
	for i := 1; i <= count; i++ {
		payload, err := payload.GeneratePayload(size)
		if err != nil {
			logger.Logger.Panicf("Error generating random string %s", err)
		}
		dataSlice = append(dataSlice, MongoData{
			ID:      i + 1,
			Payload: payload,
		})
	}
	return dataSlice
}

func WriteMongoDataToFile(filename string, data []MongoData) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}
