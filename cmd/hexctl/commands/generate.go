package commands

import (
	"os"
	"strings"

	"github.com/hanapedia/hexagon/internal/hexctl/manifest/generator"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate Kubernetes manifests for given service unit configuration.",
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

		generator.GenerateFromDirectory(inputPath, outputPath)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringVarP(&inputPath, "file", "f", "", "YAML file or directory to validate")
	generateCmd.PersistentFlags().StringVarP(&outputPath, "out", "o", "", "output directory for generated files")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
