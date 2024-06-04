package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hanapedia/hexagon/pkg/operator/object/factory"
	"github.com/hanapedia/hexagon/pkg/operator/yaml"
)

func TestDeploymentFactory(t *testing.T) {
	args := factory.DeploymentArgs{
		Name:         "test",
		Namespace:    "test",
		Image:        "test",
		Replicas:     1,
		Resource:     nil,
		Ports:        map[string]int32{"http": 8080},
		VolumeMounts: map[string]string{"config": "/config"},
		ConfigVolume: &factory.ConfigMapVolumeArgs{
			Name:  "config",
			Items: map[string]string{"config": "config.txt"},
		},
	}
	deployment := factory.NewDeployment(&args)

	// Generate the YAML
	deploymentYAML := yaml.GenerateManifest(deployment)
	// Check if some of the fields are correctly present in the generated YAML
	if !strings.Contains(string(deploymentYAML), "name: test") {
		t.Errorf("The 'name' field is incorrect or missing in the generated YAML")
	}

	if !strings.Contains(string(deploymentYAML), "namespace: test") {
		t.Errorf("The 'namespace' field is incorrect or missing in the generated YAML")
	}

	if !strings.Contains(string(deploymentYAML), "replicas: 1") {
		t.Errorf("The 'replicas' field is incorrect or missing in the generated YAML")
	}

	fmt.Printf("%s", string(deploymentYAML))
}
