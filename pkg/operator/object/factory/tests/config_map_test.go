package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hanapedia/hexagon/pkg/operator/object/factory"
	"github.com/hanapedia/hexagon/pkg/operator/yaml"
)

func TestConfigMapFactory(t *testing.T) {
	rawYaml := `---
name: gateway
ingressAdapters:
- stateless:
    action: read
    variant: rest
    route: get
  steps:
  - egressAdapter:
      stateless:
        variant: rest
        service: chain1
        action: read
        route: get`

	args := factory.ConfigMapArgs{
		Name:                   "test",
		Namespace:              "test",
		Data: map[string]string{
			"data": "test",
			"yaml": rawYaml,
		},
	}
	configMap := factory.NewConfigMap(&args)

	// Generate the YAML
	configMapYAML := yaml.GenerateManifest(configMap)
	// Check if some of the fields are correctly present in the generated YAML
	if !strings.Contains(string(configMapYAML), "name: test") {
		t.Errorf("The 'name' field is incorrect or missing in the generated YAML")
	}

	if !strings.Contains(string(configMapYAML), "namespace: test") {
		t.Errorf("The 'namespace' field is incorrect or missing in the generated YAML")
	}

	fmt.Printf("%s", string(configMapYAML))
}

