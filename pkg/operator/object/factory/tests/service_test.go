package tests

import (
	"strings"
	"testing"

	"github.com/hanapedia/hexagon/pkg/operator/object/factory"
	"github.com/hanapedia/hexagon/pkg/operator/yaml"
)

func TestServiceFactory(t *testing.T) {
	args := factory.ServiceArgs{
		Name:                   "test",
		Namespace:              "test",
		Ports:                  map[string]int32{"http": 8080},
	}
	service := factory.NewSerivce(&args)

	// Generate the YAML
	serviceYAML := yaml.GenerateManifest(service)
	// Check if some of the fields are correctly present in the generated YAML
	if !strings.Contains(string(serviceYAML), "name: test") {
		t.Errorf("The 'name' field is incorrect or missing in the generated YAML")
	}

	if !strings.Contains(string(serviceYAML), "namespace: test") {
		t.Errorf("The 'namespace' field is incorrect or missing in the generated YAML")
	}
}
