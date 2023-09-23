package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/hanapedia/the-bench/internal/tbctl/manifest/generator"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate Kubernetes manifests for given service unit configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(inputPath) == "" {
			fmt.Println("Error: Missing -f flag or empty file path")
			return
		}

		fileInfo, err := os.Stat(inputPath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if fileInfo.IsDir() {
			generator.GenerateFromDirectory(inputPath, outputPath)
		} else {
			manifestGenerator := generator.NewManifestGenerator(inputPath, outputPath)
			manifestGenerator.GenerateFromFile()
		}
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
