package loader

import (
	"fmt"
	"io/fs"
	"path/filepath"

	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/loader"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/operator/yaml"
)

func newConfigLoader(path string) loader.ConfigLoader {
	return yaml.YamlConfigLoader{Path: path}
}

func GetConfig(path string) model.ServiceUnitConfig {
	configLoader := newConfigLoader(path)
	config, err := configLoader.Load()
	if err != nil {
		logger.Logger.Fatalf("Error loading config: %v", err)
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

		if !dirEntry.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
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
