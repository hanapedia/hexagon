package loader

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hanapedia/hexagon/internal/config/application/ports"
	"github.com/hanapedia/hexagon/internal/config/infrastructure/yaml"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func newConfigLoader(path string) ports.ConfigLoader {
	return &yaml.YamlConfigLoader{Path: path}
}

func GetServiceUnitConfig(path string) *model.ServiceUnitConfig {
	configLoader := newConfigLoader(path)
	config, err := configLoader.Load()
	if err != nil {
		logger.Logger.Fatalf("Error loading config: %v", err)
	}
	return config
}

func GetClusterConfig(path string) *model.ClusterConfig {
	configLoader := newConfigLoader(path)
	config, err := configLoader.LoadClusterConfig()
	if err != nil {
		logger.Logger.Fatalf("Error loading cluster config: %v", err)
	}
	return config
}

// GetYAMLFiles retrieves all the yaml files under given directory recursively
func GetYAMLFiles(dir string) ([]string, error) {
	yamlFiles := []string{}

	err := filepath.WalkDir(dir, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}

		if dirEntry.IsDir() {
			return nil
		}

		if filepath.Base(path) == constants.CLUSTER_CONFIG_FILE_NAME {
			return nil
		}

		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			yamlFiles = append(yamlFiles, path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory %s: %v\n", dir, err)
		return nil, err
	}

	return yamlFiles, nil
}

// ReplaceInputDirectory replaces the current parent directory with new parent directory
// while preserving the other children directories
func ReplaceInputDirectory(path string, oldParent, newParent string) string {
	oldParent = filepath.Clean(oldParent)
	newParent = filepath.Clean(newParent)
	// Check if the path starts with the old parent directory
	if strings.HasPrefix(path, oldParent) {
		// Replace the old parent directory with the new parent directory
		relativePath, err := filepath.Rel(oldParent, path)
		if err != nil {
			fmt.Printf("Error getting relative path for %s: %v\n", path, err)
			return path
		}
		return filepath.Join(newParent, relativePath)
	}
	return path
}
