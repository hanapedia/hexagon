package commands

import (
	"os"
	"strings"

	"github.com/hanapedia/hexagon/internal/hexctl/adapto"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/spf13/cobra"
)

// adaptoCmd represents the generate command
var adaptoCmd = &cobra.Command{
	Use:   "adapto",
	Short: "adds adapto configuration to given service unit configurations.",
	Run: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(inputPath) == "" {
			logger.Logger.Fatalln("Missing -f flag or empty file path")
		}

		fileInfo, err := os.Stat(inputPath)
		if err != nil {
			logger.Logger.Fatalln(err)
		}
		if !fileInfo.IsDir() {
			logger.Logger.Fatalln("The input path is not a directory.")
		}

		adapto.GenerateFromDirectory(inputPath, outputPath)
	},
}

func init() {
	rootCmd.AddCommand(adaptoCmd)
	adaptoCmd.PersistentFlags().StringVarP(&inputPath, "file", "f", "", "YAML file or directory to validate")
	adaptoCmd.PersistentFlags().StringVarP(&outputPath, "out", "o", "", "output directory for generated files")
}
