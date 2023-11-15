package commands

import (
	"fmt"

	"github.com/hanapedia/hexagon/internal/datagen/mongo"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/spf13/cobra"
)

var mongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "Generate dummy data for mongo.",
	Run: func(cmd *cobra.Command, args []string) {
		size := constants.PayloadSizeMap[constants.LARGE]
		// generate large data
		data := mongo.GenerateMongoData(constants.NumInitialEntries, size)
		err := mongo.WriteMongoDataToFile(fmt.Sprintf("%s.json", constants.LARGE), data)
		if err != nil {
			logger.Logger.Panicf("Error writing to file: %s", err)
		}

		// generate medium data
		size = constants.PayloadSizeMap[constants.MEDIUM]
		data = mongo.GenerateMongoData(constants.NumInitialEntries, size)
		err = mongo.WriteMongoDataToFile(fmt.Sprintf("%s.json", constants.MEDIUM), data)
		if err != nil {
			logger.Logger.Panicf("Error writing to file: %s", err)
		}

		// generate small data
		size = constants.PayloadSizeMap[constants.SMALL]
		data = mongo.GenerateMongoData(constants.NumInitialEntries, size)
		err = mongo.WriteMongoDataToFile(fmt.Sprintf("%s.json", constants.SMALL), data)
		if err != nil {
			logger.Logger.Panicf("Error writing to file: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mongoCmd)
}
