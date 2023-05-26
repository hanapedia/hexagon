package generate

import (
	"fmt"
	"io/ioutil"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/tbctl/pkg/kubernetes/templates"
)

// take the path to the
func GenerateConfigManifest(dir string, serviceUnitConfig model.ServiceUnitConfig, path string) error {
	configData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	config := templates.ConfigConfigMap{
		Name: serviceUnitConfig.Name,
		Namespace: NAMESPACE,
		Config:    string(configData),
	}

	err = RenderAndSave(
		dir,
		fmt.Sprintf("%s-config-config-map", serviceUnitConfig.Name),
		templates.ConfigConfigMapTemplate,
		config,
	)
	return err
}

// take the path to the
func GenerateEnvManifest(dir string, serviceUnitConfig model.ServiceUnitConfig, envData string) error {
	config := templates.EnvConfigMap{
		Name: serviceUnitConfig.Name,
		Namespace: NAMESPACE,
		Envs:      envData,
	}

	err := RenderAndSave(
		dir,
		fmt.Sprintf("%s-env-config-map", serviceUnitConfig.Name),
		templates.EnvConfigMapTemplate,
		config,
	)
	return err
}
