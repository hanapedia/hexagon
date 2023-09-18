package tests

import (
	"log"
	"testing"

	// "github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/pkg/operator/yaml"
)

func TestYamlConfigLoader_InvalidIngressAdapter(t *testing.T) {
	yamlLoader := yaml.YamlConfigLoader{Path: "./testdata/invalid/invalidIngressAdapter.yaml"}
	config, err := yamlLoader.Load()
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Println(config)

}
