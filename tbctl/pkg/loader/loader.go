package loader

import (
	"fmt"
	"io/fs"
	"path/filepath"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/loader"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/yaml"
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
