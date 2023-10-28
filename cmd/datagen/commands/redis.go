package commands

import (
	"fmt"

	"github.com/hanapedia/hexagon/internal/datagen/redis"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/spf13/cobra"
)

var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Generate dummy data for redis.",
	Run: func(cmd *cobra.Command, args []string) {
	size := constants.LARGE
		// generate large data
		data := redis.GenerateRedisData(constants.NumInitialEntries, size)
		err := redis.WriteRedisDataToFile(fmt.Sprintf("%s.txt", size), data)
		if err != nil {
			logger.Logger.Panicf("Error writing to file: %s", err)
		}

		// generate medium data
		size = constants.MEDIUM
		data = redis.GenerateRedisData(constants.NumInitialEntries, size)
		err = redis.WriteRedisDataToFile(fmt.Sprintf("%s.txt", size), data)
		if err != nil {
			logger.Logger.Panicf("Error writing to file: %s", err)
		}

		// generate small data
		size = constants.SMALL
		data = redis.GenerateRedisData(constants.NumInitialEntries, size)
		err = redis.WriteRedisDataToFile(fmt.Sprintf("%s.txt", size), data)
		if err != nil {
			logger.Logger.Panicf("Error writing to file: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(redisCmd)
}

