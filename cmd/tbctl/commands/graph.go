package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/hanapedia/the-bench/internal/tbctl/graph"
	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Generate GraphML representation of service unit application.",
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
			graph.GenerateFromDirectory(inputPath, outputPath)
		} else {
			fmt.Println("Error: Only genetation from directory supported. Please provide directory path.")
			return
		}
},
}

func init() {
	rootCmd.AddCommand(graphCmd)
	rootCmd.PersistentFlags().StringVarP(&inputPath, "file", "f", "", "YAML file or directory to validate")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "out", "o", "./", "output directory for generated files")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// graphCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// graphCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
