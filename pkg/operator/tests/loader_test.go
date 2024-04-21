package tests

import (
	"log"
	"testing"

	// "github.com/hanapedia/hexagon/config/model"
	"github.com/hanapedia/hexagon/internal/config/infrastructure/yaml"
)

func TestYamlConfigLoader_InvalidIngressAdapter(t *testing.T) {
	yamlLoader := yaml.YamlConfigLoader{Path: "./testdata/invalid/invalidIngressAdapter.yaml"}
	config, err := yamlLoader.Load()
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Println(config)

}
