/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hanapedia/the-bench/tbctl/pkg/validation"
	"github.com/spf13/cobra"
)

var filePath string

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate service unit configs from a YAML file or directory containing YAML files.",
	Run: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(filePath) == "" {
			fmt.Println("Error: Missing -f flag or empty file path")
			return
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if fileInfo.IsDir() {
			err := validation.ValidateDirectory(filePath)
			if len(err.ServiceUnitFieldErrors) > 0 || len(err.AdapterFieldErrors) > 0 || len(err.MappingErrors) > 0 {
				log.Fatalf(
					"Validation failed with %v service unit field errors, %v adapter field errors, and %v mapping errors.",
					len(err.ServiceUnitFieldErrors),
					len(err.AdapterFieldErrors),
					len(err.MappingErrors),
				)
			}
		} else {
			sufe, aef := validation.ValidateFile(filePath)
			if len(sufe) > 0 || len(aef) > 0 {
				log.Fatalf("Validation failed with %v service unit field errors and %v adapter field errors.", len(sufe), len(aef))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "YAML file or directory to validate")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
